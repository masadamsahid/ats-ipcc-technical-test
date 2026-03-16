package repository

import (
	"testing"
	"tsilodot/model"
)

func TestUserRepository_CreateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	createdUser, err := repo.CreateUser(nil, user)
	if err != nil {
		t.Errorf("CreateUser() error = %v", err)
		return
	}

	if createdUser.ID == 0 {
		t.Errorf("CreateUser() ID = 0, want non-zero")
	}
	if createdUser.Email != "test@example.com" {
		t.Errorf("CreateUser() Email = %v, want %v", createdUser.Email, "test@example.com")
	}
}

func TestUserRepository_FindUserByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{Name: "Test User", Email: "test@example.com", Password: "password123"}
	createdUser, _ := repo.CreateUser(nil, user)

	foundUser, err := repo.FindUserByID(nil, createdUser.ID)
	if err != nil {
		t.Errorf("FindUserByID() error = %v", err)
		return
	}

	if foundUser.ID != createdUser.ID {
		t.Errorf("FindUserByID() ID = %v, want %v", foundUser.ID, createdUser.ID)
	}
}

func TestUserRepository_FindUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{Name: "Test User", Email: "test@example.com", Password: "password123"}
	repo.CreateUser(nil, user)

	foundUser, err := repo.FindUserByEmail(nil, "test@example.com")
	if err != nil {
		t.Errorf("FindUserByEmail() error = %v", err)
		return
	}

	if foundUser.Email != "test@example.com" {
		t.Errorf("FindUserByEmail() Email = %v, want %v", foundUser.Email, "test@example.com")
	}
}
