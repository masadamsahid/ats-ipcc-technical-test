package dto

import "time"

type TaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"required,oneof=pending completed"`
	DueDate     string `json:"due_date" validate:"required,datetime=2006-01-02"`
}

func (r *TaskRequest) GetDueDate() *time.Time {
	if r.DueDate == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", r.DueDate)
	if err != nil {
		return nil
	}
	return &t
}

type TaskResponse struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
