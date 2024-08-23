package dao

import (
	"errors"

	"github.com/0b0e0e7c/chat/model"
	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// FindUserByUsernameAndPassword 根据用户名和密码查找用户
func (dao *UserDAO) FindUserByUsernameAndPassword(username, password string) (loginUser *model.User, err error) {
	result := dao.db.Where("username = ? AND password = ?", username, password).First(&loginUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, result.Error
	}
	return
}

// CreateUserByUsernameAndPassword 创建用户
func (dao *UserDAO) CreateUserByUsernameAndPassword(username, password string) (newUser *model.User, err error) {
	newUser = &model.User{
		Username: username,
		Password: password,
	}
	result := dao.db.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return
}
