package repository

import (
	"testing"
	"tsilodot/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&model.User{}, &model.Task{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestTaskRepository_CreateTask(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	task := &model.Task{
		UserID:      1,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
	}

	createdTask, err := repo.CreateTask(nil, task)
	if err != nil {
		t.Errorf("CreateTask() error = %v", err)
		return
	}

	if createdTask.ID == 0 {
		t.Errorf("CreateTask() ID = 0, want non-zero")
	}
	if createdTask.Title != "Test Task" {
		t.Errorf("CreateTask() Title = %v, want %v", createdTask.Title, "Test Task")
	}
}

func TestTaskRepository_FindTasksByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	userId := uint(1)
	repo.CreateTask(nil, &model.Task{UserID: userId, Title: "Task 1", Status: "pending"})
	repo.CreateTask(nil, &model.Task{UserID: userId, Title: "Task 2", Status: "completed"})
	repo.CreateTask(nil, &model.Task{UserID: 2, Title: "Task 3", Status: "pending"})

	tasks, total, err := repo.FindTasksByUserID(nil, userId, 5, 0)
	if err != nil {
		t.Errorf("FindTasksByUserID() error = %v", err)
		return
	}

	if total != 2 {
		t.Errorf("FindTasksByUserID() total = %v, want 2", total)
	}
	if len(tasks) != 2 {
		t.Errorf("FindTasksByUserID() len(tasks) = %v, want 2", len(tasks))
	}
}

func TestTaskRepository_FindTaskByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	task := &model.Task{UserID: 1, Title: "Test Task", Status: "pending"}
	createdTask, _ := repo.CreateTask(nil, task)

	foundTask, err := repo.FindTaskByID(nil, createdTask.ID)
	if err != nil {
		t.Errorf("FindTaskByID() error = %v", err)
		return
	}

	if foundTask.ID != createdTask.ID {
		t.Errorf("FindTaskByID() ID = %v, want %v", foundTask.ID, createdTask.ID)
	}
}

func TestTaskRepository_UpdateTask(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	task := &model.Task{UserID: 1, Title: "Old Title", Status: "pending"}
	createdTask, _ := repo.CreateTask(nil, task)

	createdTask.Title = "New Title"
	updatedTask, err := repo.UpdateTask(nil, createdTask)
	if err != nil {
		t.Errorf("UpdateTask() error = %v", err)
		return
	}

	if updatedTask.Title != "New Title" {
		t.Errorf("UpdateTask() Title = %v, want %v", updatedTask.Title, "New Title")
	}
}

func TestTaskRepository_DeleteTask(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	task := &model.Task{UserID: 1, Title: "To Delete", Status: "pending"}
	createdTask, _ := repo.CreateTask(nil, task)

	err := repo.DeleteTask(nil, createdTask.ID)
	if err != nil {
		t.Errorf("DeleteTask() error = %v", err)
		return
	}

	_, err = repo.FindTaskByID(nil, createdTask.ID)
	if err == nil {
		t.Errorf("FindTaskByID() error = nil, want error (not found)")
	}
}
