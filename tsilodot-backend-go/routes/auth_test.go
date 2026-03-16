package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"tsilodot/controller"
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

func TestSetupAuthRoutes(t *testing.T) {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{
			validate: validator.New(),
		},
	})
	api := app.Group("/api")

	mockSvc := &MockAuthService{}
	authController := controller.NewAuthController(mockSvc)

	SetupAuthRoutes(api, authController)

	t.Run("POST /api/auth/register", func(t *testing.T) {
		mockSvc.RegisterFunc = func(params *model.User) (*dto.RegisterResponseData, error) {
			return &dto.RegisterResponseData{ID: 1, Email: params.Email}, nil
		}

		reqBody, _ := json.Marshal(dto.RegisterRequest{
			Name:            "Test",
			Email:           "test@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		})
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/auth/login", func(t *testing.T) {
		mockSvc.LoginFunc = func(params *dto.LoginRequest) (*dto.LoginResponseData, error) {
			return &dto.LoginResponseData{AccessToken: "token"}, nil
		}

		reqBody, _ := json.Marshal(dto.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		})
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}
