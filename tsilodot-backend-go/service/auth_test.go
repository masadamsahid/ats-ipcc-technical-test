package service

import (
	"errors"
	"testing"
	"tsilodot/dto"
	"tsilodot/helpers"
	"tsilodot/model"

	"gorm.io/gorm"
)

type MockUserRepository struct {
	users map[uint]*model.User
	email map[string]*model.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint]*model.User),
		email: make(map[string]*model.User),
	}
}

func (m *MockUserRepository) CreateUser(db *gorm.DB, param *model.User) (*model.User, error) {
	param.ID = uint(len(m.users) + 1)
	m.users[param.ID] = param
	m.email[param.Email] = param
	return param, nil
}

func (m *MockUserRepository) FindUserByID(db *gorm.DB, userId uint) (*model.User, error) {
	user, ok := m.users[userId]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) FindUserByEmail(db *gorm.DB, email string) (*model.User, error) {
	user, ok := m.email[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func TestAuthService_Register(t *testing.T) {
	repo := NewMockUserRepository()
	svc := NewAuthService(repo)

	user := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	res, err := svc.Register(user)
	if err != nil {
		t.Errorf("Register() error = %v", err)
		return
	}

	if res.Name != "Test User" {
		t.Errorf("Register() Name = %v, want %v", res.Name, "Test User")
	}
	if res.AccessToken == "" {
		t.Errorf("Register() AccessToken is empty")
	}
}

func TestAuthService_Login(t *testing.T) {
	repo := NewMockUserRepository()
	svc := NewAuthService(repo)

	// Register a user first
	pwd := "password123"
	hashedPwd, _ := helpers.HashPassword(pwd)
	repo.CreateUser(nil, &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPwd,
	})

	loginReq := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: pwd,
	}

	res, err := svc.Login(loginReq)
	if err != nil {
		t.Errorf("Login() error = %v", err)
		return
	}

	if res.Email != "test@example.com" {
		t.Errorf("Login() Email = %v, want %v", res.Email, "test@example.com")
	}
	if res.AccessToken == "" {
		t.Errorf("Login() AccessToken is empty")
	}
}
