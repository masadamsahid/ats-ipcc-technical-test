package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"tsilodot/dto"
	"tsilodot/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

type MockAuthService struct {
	RegisterFunc func(params *model.User) (*dto.RegisterResponseData, error)
	LoginFunc    func(params *dto.LoginRequest) (*dto.LoginResponseData, error)
}

func (m *MockAuthService) Register(params *model.User) (*dto.RegisterResponseData, error) {
	return m.RegisterFunc(params)
}

func (m *MockAuthService) Login(params *dto.LoginRequest) (*dto.LoginResponseData, error) {
	return m.LoginFunc(params)
}

func TestAuthController_Register(t *testing.T) {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{
			validate: validator.New(),
		},
	})
	mockSvc := &MockAuthService{}
	ctrl := NewAuthController(mockSvc)

	app.Post("/register", ctrl.Register)

	t.Run("Success", func(t *testing.T) {
		mockSvc.RegisterFunc = func(params *model.User) (*dto.RegisterResponseData, error) {
			return &dto.RegisterResponseData{
				ID:    1,
				Name:  params.Name,
				Email: params.Email,
			}, nil
		}

		reqBody, _ := json.Marshal(dto.RegisterRequest{
			Name:            "Test User",
			Email:           "test@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		})
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", resp.StatusCode)
		}

		var response dto.ResponseGeneric[dto.RegisterResponseData]
		json.NewDecoder(resp.Body).Decode(&response)

		if response.Message != "Registration successful" {
			t.Errorf("Expected message 'Registration successful', got '%s'", response.Message)
		}
		if response.Data == nil || response.Data.Email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%v'", response.Data)
		}
	})

	t.Run("Validation Error", func(t *testing.T) {
		reqBody, _ := json.Marshal(dto.RegisterRequest{
			Name:     "",
			Email:    "invalid-email",
			Password: "short",
		})
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var response dto.ResponseGeneric[any]
		json.NewDecoder(resp.Body).Decode(&response)
		if response.Message != "Failed" {
			t.Errorf("Expected message 'Failed', got '%s'", response.Message)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc.RegisterFunc = func(params *model.User) (*dto.RegisterResponseData, error) {
			return nil, errors.New("db error")
		}

		reqBody, _ := json.Marshal(dto.RegisterRequest{
			Name:            "Test User",
			Email:           "test@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		})
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	})
}

func TestAuthController_Login(t *testing.T) {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{
			validate: validator.New(),
		},
	})
	mockSvc := &MockAuthService{}
	ctrl := NewAuthController(mockSvc)

	app.Post("/login", ctrl.Login)

	t.Run("Success", func(t *testing.T) {
		mockSvc.LoginFunc = func(params *dto.LoginRequest) (*dto.LoginResponseData, error) {
			return &dto.LoginResponseData{
				ID:          1,
				Email:       params.Email,
				AccessToken: "fake-token",
			}, nil
		}

		reqBody, _ := json.Marshal(dto.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		})
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response dto.ResponseGeneric[dto.LoginResponseData]
		json.NewDecoder(resp.Body).Decode(&response)

		if response.Data == nil || response.Data.AccessToken != "fake-token" {
			t.Errorf("Expected token 'fake-token', got '%v'", response.Data)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockSvc.LoginFunc = func(params *dto.LoginRequest) (*dto.LoginResponseData, error) {
			return nil, errors.New("unauthorized")
		}

		reqBody, _ := json.Marshal(dto.LoginRequest{
			Email:    "wrong@example.com",
			Password: "wrongpassword",
		})
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}
