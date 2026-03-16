package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"tsilodot/dto"
	"tsilodot/helpers"
	"tsilodot/model"
	"tsilodot/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type MockTaskService struct {
	CreateTaskFunc       func(userId uint, task *model.Task) (*model.Task, error)
	GetTasksByUserIDFunc func(userId uint, limit int, offset int) ([]model.Task, int64, error)
	GetTaskByIDFunc      func(taskId uint, userId uint) (*model.Task, error)
	UpdateTaskFunc       func(taskId uint, userId uint, task *model.Task) (*model.Task, error)
	DeleteTaskFunc       func(taskId uint, userId uint) error
}

func (m *MockTaskService) CreateTask(userId uint, task *model.Task) (*model.Task, error) {
	return m.CreateTaskFunc(userId, task)
}

func (m *MockTaskService) GetTasksByUserID(userId uint, limit int, offset int) ([]model.Task, int64, error) {
	return m.GetTasksByUserIDFunc(userId, limit, offset)
}

func (m *MockTaskService) GetTaskByID(taskId uint, userId uint) (*model.Task, error) {
	return m.GetTaskByIDFunc(taskId, userId)
}

func (m *MockTaskService) UpdateTask(taskId uint, userId uint, task *model.Task) (*model.Task, error) {
	return m.UpdateTaskFunc(taskId, userId, task)
}

func (m *MockTaskService) DeleteTask(taskId uint, userId uint) error {
	return m.DeleteTaskFunc(taskId, userId)
}

func setupTaskApp(svc service.ITaskService) *fiber.App {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{
			validate: validator.New(),
		},
	})
	ctrl := NewTaskController(svc)

	// Middleware to inject user claims
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user", &helpers.AuthTokenClaims{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		})
		return c.Next()
	})

	app.Post("/tasks", ctrl.CreateTask)
	app.Get("/tasks", ctrl.GetTasks)
	app.Get("/tasks/:id", ctrl.GetTaskByID)
	app.Put("/tasks/:id", ctrl.UpdateTask)
	app.Delete("/tasks/:id", ctrl.DeleteTask)

	return app
}

func TestTaskController_CreateTask(t *testing.T) {
	mockSvc := &MockTaskService{}
	app := setupTaskApp(mockSvc)

	t.Run("Success", func(t *testing.T) {
		mockSvc.CreateTaskFunc = func(userId uint, task *model.Task) (*model.Task, error) {
			task.ID = 1
			return task, nil
		}

		reqBody, _ := json.Marshal(dto.TaskRequest{
			Title:       "Test Task",
			Description: "Test Description",
			Status:      "pending",
			DueDate:     "2026-03-20",
		})
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response dto.ResponseGeneric[model.Task]
		json.NewDecoder(resp.Body).Decode(&response)

		if response.Data == nil || response.Data.Title != "Test Task" {
			t.Errorf("Expected title 'Test Task', got '%v'", response.Data)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc.CreateTaskFunc = func(userId uint, task *model.Task) (*model.Task, error) {
			return nil, errors.New("db error")
		}

		reqBody, _ := json.Marshal(dto.TaskRequest{
			Title:   "Test Task",
			Status:  "pending",
			DueDate: "2026-03-20",
		})
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	})
}

func TestTaskController_GetTasks(t *testing.T) {
	mockSvc := &MockTaskService{}
	app := setupTaskApp(mockSvc)

	t.Run("Success", func(t *testing.T) {
		mockSvc.GetTasksByUserIDFunc = func(userId uint, limit int, offset int) ([]model.Task, int64, error) {
			return []model.Task{{Title: "Task 1"}, {Title: "Task 2"}}, 2, nil
		}

		req, _ := http.NewRequest("GET", "/tasks?page=1&limit=5", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response dto.ResponseGeneric[[]model.Task]
		json.NewDecoder(resp.Body).Decode(&response)

		if response.Data == nil || len(*response.Data) != 2 {
			t.Errorf("Expected 2 tasks, got %v", response.Data)
		}
		if response.Pagination.TotalItems != 2 {
			t.Errorf("Expected total items 2, got %d", response.Pagination.TotalItems)
		}
	})
}

func TestTaskController_GetTaskByID(t *testing.T) {
	mockSvc := &MockTaskService{}
	app := setupTaskApp(mockSvc)

	t.Run("Success", func(t *testing.T) {
		mockSvc.GetTaskByIDFunc = func(taskId uint, userId uint) (*model.Task, error) {
			return &model.Task{ID: taskId, Title: "Task 1"}, nil
		}

		req, _ := http.NewRequest("GET", "/tasks/1", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response dto.ResponseGeneric[model.Task]
		json.NewDecoder(resp.Body).Decode(&response)

		if response.Data == nil || response.Data.ID != 1 {
			t.Errorf("Expected task ID 1, got %v", response.Data)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		mockSvc.GetTaskByIDFunc = func(taskId uint, userId uint) (*model.Task, error) {
			return nil, errors.New("not found")
		}

		req, _ := http.NewRequest("GET", "/tasks/999", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tasks/abc", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestTaskController_UpdateTask(t *testing.T) {
	mockSvc := &MockTaskService{}
	app := setupTaskApp(mockSvc)

	t.Run("Success", func(t *testing.T) {
		mockSvc.UpdateTaskFunc = func(taskId uint, userId uint, task *model.Task) (*model.Task, error) {
			task.ID = taskId
			return task, nil
		}

		reqBody, _ := json.Marshal(dto.TaskRequest{
			Title:   "Updated Task",
			Status:  "completed",
			DueDate: "2026-03-21",
		})
		req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response dto.ResponseGeneric[model.Task]
		json.NewDecoder(resp.Body).Decode(&response)

		if response.Data == nil || response.Data.Title != "Updated Task" {
			t.Errorf("Expected title 'Updated Task', got '%v'", response.Data)
		}
	})
}

func TestTaskController_DeleteTask(t *testing.T) {
	mockSvc := &MockTaskService{}
	app := setupTaskApp(mockSvc)

	t.Run("Success", func(t *testing.T) {
		mockSvc.DeleteTaskFunc = func(taskId uint, userId uint) error {
			return nil
		}

		req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response dto.ResponseGeneric[any]
		json.NewDecoder(resp.Body).Decode(&response)

		if response.Message != "Task deleted successfully" {
			t.Errorf("Expected message 'Task deleted successfully', got '%s'", response.Message)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockSvc.DeleteTaskFunc = func(taskId uint, userId uint) error {
			return errors.New("unauthorized")
		}

		req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	})
}
