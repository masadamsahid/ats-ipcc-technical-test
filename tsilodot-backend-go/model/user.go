package model

import "time"

type User struct {
	ID       uint   `gorm:"column:id;primary_key" json:"id"`
	Name     string `gorm:"type:text;not null" json:"name"`
	Email    string `gorm:"type:text;uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:text;not null" json:"-"`

	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
