package service

import (
	"errors"
	"tsilodot/model"
	"tsilodot/repository"
)

type TaskService struct {
	taskRepository repository.ITaskRepository
}

type ITaskService interface {
	CreateTask(userId uint, task *model.Task) (*model.Task, error)
	GetTasksByUserID(userId uint, limit int, offset int) ([]model.Task, int64, error)
	GetTaskByID(taskId uint, userId uint) (*model.Task, error)
	UpdateTask(taskId uint, userId uint, task *model.Task) (*model.Task, error)
	DeleteTask(taskId uint, userId uint) error
}

func NewTaskService(taskRepository repository.ITaskRepository) ITaskService {
	return &TaskService{taskRepository: taskRepository}
}

func (s *TaskService) CreateTask(userId uint, task *model.Task) (*model.Task, error) {
	task.UserID = userId
	return s.taskRepository.CreateTask(nil, task)
}

func (s *TaskService) GetTasksByUserID(userId uint, limit int, offset int) ([]model.Task, int64, error) {
	return s.taskRepository.FindTasksByUserID(nil, userId, limit, offset)
}

func (s *TaskService) GetTaskByID(taskId uint, userId uint) (*model.Task, error) {
	task, err := s.taskRepository.FindTaskByID(nil, taskId)
	if err != nil {
		return nil, err
	}

	if task.UserID != userId {
		return nil, errors.New("unauthorized: task does not belong to user")
	}

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

	return s.taskRepository.UpdateTask(nil, task)
}

func (s *TaskService) DeleteTask(taskId uint, userId uint) error {
	task, err := s.taskRepository.FindTaskByID(nil, taskId)
	if err != nil {
		return err
	}

	if task.UserID != userId {
		return errors.New("unauthorized: cannot delete task belonging to another user")
	}

	return s.taskRepository.DeleteTask(nil, taskId)
}
