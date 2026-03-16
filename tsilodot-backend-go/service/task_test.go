package service

import (
	"errors"
	"testing"
	"tsilodot/model"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MockTaskRepository struct {
	tasks map[uint]*model.Task
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		tasks: make(map[uint]*model.Task),
	}
}

func (m *MockTaskRepository) CreateTask(db *gorm.DB, param *model.Task) (*model.Task, error) {
	param.ID = uint(len(m.tasks) + 1)
	m.tasks[param.ID] = param
	return param, nil
}

func (m *MockTaskRepository) FindTasksByUserID(db *gorm.DB, userId uint, limit int, offset int) ([]model.Task, int64, error) {
	var tasks []model.Task
	for _, t := range m.tasks {
		if t.UserID == userId {
			tasks = append(tasks, *t)
		}
	}
	return tasks, int64(len(tasks)), nil
}

func (m *MockTaskRepository) FindTaskByID(db *gorm.DB, taskId uint) (*model.Task, error) {
	task, ok := m.tasks[taskId]
	if !ok {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (m *MockTaskRepository) UpdateTask(db *gorm.DB, task *model.Task) (*model.Task, error) {
	m.tasks[task.ID] = task
	return task, nil
}

func (m *MockTaskRepository) DeleteTask(db *gorm.DB, taskId uint) error {
	delete(m.tasks, taskId)
	return nil
}

func TestTaskService_GetTaskByID_Authorization(t *testing.T) {
	s := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: s.Addr()})

	repo := NewMockTaskRepository()
	svc := NewTaskService(repo, redisClient)

	// Create task for user 1
	task, _ := repo.CreateTask(nil, &model.Task{UserID: 1, Title: "User 1 Task"})

	// Try to get with user 2
	_, err := svc.GetTaskByID(task.ID, 2)
	if err == nil {
		t.Errorf("GetTaskByID() expected error for unauthorized user, got nil")
	}

	// Get with user 1
	foundTask, err := svc.GetTaskByID(task.ID, 1)
	if err != nil {
		t.Errorf("GetTaskByID() error = %v", err)
		return
	}
	if foundTask.ID != task.ID {
		t.Errorf("GetTaskByID() ID = %v, want %v", foundTask.ID, task.ID)
	}
}

func TestTaskService_UpdateTask_Authorization(t *testing.T) {
	s := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: s.Addr()})

	repo := NewMockTaskRepository()
	svc := NewTaskService(repo, redisClient)

	task, _ := repo.CreateTask(nil, &model.Task{UserID: 1, Title: "Original Title"})

	// Try update with user 2
	_, err := svc.UpdateTask(task.ID, 2, &model.Task{Title: "New Title"})
	if err == nil {
		t.Errorf("UpdateTask() expected error for unauthorized user, got nil")
	}

	// Update with user 1
	updatedTask, err := svc.UpdateTask(task.ID, 1, &model.Task{Title: "New Title"})
	if err != nil {
		t.Errorf("UpdateTask() error = %v", err)
		return
	}
	if updatedTask.Title != "New Title" {
		t.Errorf("UpdateTask() Title = %v, want %v", updatedTask.Title, "New Title")
	}
}

func TestTaskService_DeleteTask_Authorization(t *testing.T) {
	s := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: s.Addr()})

	repo := NewMockTaskRepository()
	svc := NewTaskService(repo, redisClient)

	task, _ := repo.CreateTask(nil, &model.Task{UserID: 1, Title: "To Delete"})

	// Try delete with user 2
	err := svc.DeleteTask(task.ID, 2)
	if err == nil {
		t.Errorf("DeleteTask() expected error for unauthorized user, got nil")
	}

	// Delete with user 1
	err = svc.DeleteTask(task.ID, 1)
	if err != nil {
		t.Errorf("DeleteTask() error = %v", err)
	}
}
