package model

import "time"

type Task struct {
	ID          uint       `gorm:"column:id;primary_key" json:"id"`
	UserID      uint       `gorm:"column:user_id;not null" json:"user_id"`
	Title       string     `gorm:"type:text;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      string     `gorm:"type:text;default:'pending'" json:"status"`
	DueDate     *time.Time `gorm:"type:timestamp" json:"due_date"`

	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (Task) TableName() string {
	return "tasks"
}
