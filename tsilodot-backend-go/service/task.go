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
	"github.com/rs/zerolog/log"
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
		if errUnmarshal := json.Unmarshal([]byte(val), &task); errUnmarshal == nil {

			if task.UserID != userId {
				log.Warn().Uint("task_id", taskId).Uint("user_id", userId).Msg("Unauthorized cache access attempt")
				return nil, errors.New("unauthorized: task does not belong to user")
			}
			return &task, nil
		} else {
			log.Error().Err(errUnmarshal).Str("cache_key", cacheKey).Msg("Error unmarshaling cached task")
		}
	} else if err != redis.Nil {
		log.Error().Err(err).Str("cache_key", cacheKey).Msg("Error getting task from Redis")
	}

	task, err := s.taskRepository.FindTaskByID(nil, taskId)
	if err != nil {
		return nil, err
	}

	if task.UserID != userId {
		log.Warn().Uint("task_id", taskId).Uint("user_id", userId).Msg("Unauthorized task access attempt")
		return nil, errors.New("unauthorized: task does not belong to user")
	}

	// Cache to redis with 10 mins
	taskData, errMarshal := json.Marshal(task)
	if errMarshal != nil {
		log.Error().Err(errMarshal).Uint("task_id", taskId).Msg("Error marshaling task for cache")
	} else {
		errSet := s.redisClient.Set(context.Background(), cacheKey, taskData, 10*time.Minute).Err()
		if errSet != nil {
			log.Error().Err(errSet).Str("cache_key", cacheKey).Msg("Error setting task in Redis")
		}
	}

	return task, nil
}

func (s *TaskService) UpdateTask(taskId uint, userId uint, taskUpdate *model.Task) (*model.Task, error) {
	task, err := s.taskRepository.FindTaskByID(nil, taskId)
	if err != nil {
		return nil, err
	}

	if task.UserID != userId {
		log.Warn().Uint("task_id", taskId).Uint("user_id", userId).Msg("Unauthorized task update attempt")
		return nil, errors.New("unauthorized: cannot update task belonging to another user")
	}

	task.Title = taskUpdate.Title
	task.Description = taskUpdate.Description
	task.Status = taskUpdate.Status
	task.DueDate = taskUpdate.DueDate

	updatedTask, err := s.taskRepository.UpdateTask(nil, task)
	if err == nil { // Invalidate cache
		cacheKey := fmt.Sprintf("task:%d", taskId)
		errDel := s.redisClient.Del(context.Background(), cacheKey).Err()
		if errDel != nil {
			log.Error().Err(errDel).Str("cache_key", cacheKey).Msg("Error deleting task from Redis")
		}
	}

	return updatedTask, err
}

func (s *TaskService) DeleteTask(taskId uint, userId uint) error {
	task, err := s.taskRepository.FindTaskByID(nil, taskId)
	if err != nil {
		return err
	}

	if task.UserID != userId {
		log.Warn().Uint("task_id", taskId).Uint("user_id", userId).Msg("Unauthorized task deletion attempt")
		return errors.New("unauthorized: cannot delete task belonging to another user")
	}

	err = s.taskRepository.DeleteTask(nil, taskId)
	if err == nil { // Invalidate cache
		cacheKey := fmt.Sprintf("task:%d", taskId)
		errDel := s.redisClient.Del(context.Background(), cacheKey).Err()
		if errDel != nil {
			log.Error().Err(errDel).Str("cache_key", cacheKey).Msg("Error deleting task from Redis after deletion")
		}
	}

	return err
}
