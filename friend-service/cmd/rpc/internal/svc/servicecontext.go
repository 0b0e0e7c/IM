package svc

import (
	"github.com/0b0e0e7c/IM/friend-service/cmd/rpc/internal/config"
	"github.com/0b0e0e7c/IM/friend-service/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Friend{})

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
