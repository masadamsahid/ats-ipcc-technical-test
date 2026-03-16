package repository

import (
	"tsilodot/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	CreateUser(db *gorm.DB, param *model.User) (*model.User, error)
	FindUserByID(db *gorm.DB, userId uint) (*model.User, error)
	FindUserByEmail(db *gorm.DB, email string) (*model.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CreateUser(db *gorm.DB, param *model.User) (*model.User, error) {
	if db == nil {
		db = u.db
	}
	user := model.User{
		Name:     param.Name,
		Email:    param.Email,
		Password: param.Password,
	}

	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) FindUserByID(db *gorm.DB, userId uint) (*model.User, error) {
	if db == nil {
		db = u.db
	}

	var user model.User
	err := db.First(&user, userId).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) FindUserByEmail(db *gorm.DB, email string) (*model.User, error) {
	if db == nil {
		db = u.db
	}

	var user model.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
