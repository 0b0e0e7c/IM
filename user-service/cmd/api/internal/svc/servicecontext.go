package svc

import (
	"IM/user-service/cmd/api/internal/config"
	"IM/user-service/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

// NewServiceContext return a new ServiceContext
func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移
	db.AutoMigrate(&model.User{})

	config.DB = db

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
