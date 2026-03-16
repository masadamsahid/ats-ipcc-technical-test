package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"tsilodot/model"
	"tsilodot/repository"

	"github.com/redis/go-redis/v9"
)

type TaskService struct {
	taskRepository repository.ITaskRepository
	redisClient    *redis.Client
}

type ITaskService interface {
	CreateTask(userId uint, task *model.Task) (*model.Task, error)
	GetTasksByUserID(userId uint, limit int, offset int) ([]model.Task, int64, error)
	GetTaskByID(taskId uint, userId uint) (*model.Task, error)
	UpdateTask(taskId uint, userId uint, task *model.Task) (*model.Task, error)
	DeleteTask(taskId uint, userId uint) error
}

func NewTaskService(taskRepository repository.ITaskRepository, redisClient *redis.Client) ITaskService {
	return &TaskService{
		taskRepository: taskRepository,
		redisClient:    redisClient,
	}
}

func (s *TaskService) CreateTask(userId uint, task *model.Task) (*model.Task, error) {
	task.UserID = userId
	return s.taskRepository.CreateTask(nil, task)
}

func (s *TaskService) GetTasksByUserID(userId uint, limit int, offset int) ([]model.Task, int64, error) {
	return s.taskRepository.FindTasksByUserID(nil, userId, limit, offset)
}

func (s *TaskService) GetTaskByID(taskId uint, userId uint) (*model.Task, error) {
	cacheKey := fmt.Sprintf("task:%d", taskId)

	// get from Redis first
	val, err := s.redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var task model.Task
		if json.Unmarshal([]byte(val), &task) == nil {

			if task.UserID != userId {
				return nil, errors.New("unauthorized: task does not belong to user")
			}
			return &task, nil
		}
	}

	task, err := s.taskRepository.FindTaskByID(nil, taskId)
	if err != nil {
		return nil, err
	}

	if task.UserID != userId {
		return nil, errors.New("unauthorized: task does not belong to user")
	}

	// Cache to redis with 10 mins
	taskData, _ := json.Marshal(task)
	s.redisClient.Set(context.Background(), cacheKey, taskData, 10*time.Minute)

	return task, nil
}

func (s *TaskService) UpdateTask(taskId uint, userId uint, taskUpdate *model.Task) (*model.Task, error) {
	task, err := s.taskRepository.FindTaskByID(nil, taskId)
	if err != nil {
		return nil, err
	}

	if task.UserID != userId {
		return nil, errors.New("unauthorized: cannot update task belonging to another user")
	}

	task.Title = taskUpdate.Title
	task.Description = taskUpdate.Description
	task.Status = taskUpdate.Status
	task.DueDate = taskUpdate.DueDate

	updatedTask, err := s.taskRepository.UpdateTask(nil, task)
	if err == nil { // Invalidate cache
		cacheKey := fmt.Sprintf("task:%d", taskId)
		s.redisClient.Del(context.Background(), cacheKey)
	}

	return updatedTask, err
}

func (s *TaskService) DeleteTask(taskId uint, userId uint) error {
	task, err := s.taskRepository.FindTaskByID(nil, taskId)
	if err != nil {
		return err
	}

	if task.UserID != userId {
		return errors.New("unauthorized: cannot delete task belonging to another user")
	}

	err = s.taskRepository.DeleteTask(nil, taskId)
	if err == nil { // Invalidate cache
		cacheKey := fmt.Sprintf("task:%d", taskId)
		s.redisClient.Del(context.Background(), cacheKey)
	}

	return err
}
