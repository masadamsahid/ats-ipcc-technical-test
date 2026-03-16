package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"tsilodot/controller"
	"tsilodot/dto"
	"tsilodot/helpers"
	"tsilodot/model"

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

func TestSetupTaskRoutes(t *testing.T) {
	helpers.InitJWT()
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{
			validate: validator.New(),
		},
	})
	api := app.Group("/api")

	mockSvc := &MockTaskService{}
	taskController := controller.NewTaskController(mockSvc)

	SetupTaskRoutes(api, taskController)

	token, _ := helpers.CreateAuthToken(helpers.AuthTokenClaims{
		ID:    1,
		Email: "test@example.com",
	})
	authHeader := fmt.Sprintf("Bearer %s", token)

	t.Run("GET /api/tasks (Protected - Success)", func(t *testing.T) {
		mockSvc.GetTasksByUserIDFunc = func(userId uint, limit int, offset int) ([]model.Task, int64, error) {
			return []model.Task{}, 0, nil
		}

		req, _ := http.NewRequest("GET", "/api/tasks", nil)
		req.Header.Set("Authorization", authHeader)

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/tasks (Protected - Unauthorized)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/tasks", nil)

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/tasks", func(t *testing.T) {
		mockSvc.CreateTaskFunc = func(userId uint, task *model.Task) (*model.Task, error) {
			return task, nil
		}

		reqBody, _ := json.Marshal(dto.TaskRequest{
			Title:   "Test Task",
			Status:  "pending",
			DueDate: "2026-03-16",
		})
		req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/tasks/:id", func(t *testing.T) {
		mockSvc.GetTaskByIDFunc = func(taskId uint, userId uint) (*model.Task, error) {
			return &model.Task{Title: "Test Task"}, nil
		}

		req, _ := http.NewRequest("GET", "/api/tasks/1", nil)
		req.Header.Set("Authorization", authHeader)

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/tasks/:id", func(t *testing.T) {
		mockSvc.UpdateTaskFunc = func(taskId uint, userId uint, task *model.Task) (*model.Task, error) {
			return task, nil
		}

		reqBody, _ := json.Marshal(dto.TaskRequest{
			Title:   "Updated Task",
			Status:  "completed",
			DueDate: "2026-03-17",
		})
		req, _ := http.NewRequest("PUT", "/api/tasks/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/tasks/:id", func(t *testing.T) {
		mockSvc.DeleteTaskFunc = func(taskId uint, userId uint) error {
			return nil
		}

		req, _ := http.NewRequest("DELETE", "/api/tasks/1", nil)
		req.Header.Set("Authorization", authHeader)

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}
