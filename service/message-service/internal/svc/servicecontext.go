package svc

import (
	"github.com/0b0e0e7c/IM/model"
	"github.com/0b0e0e7c/IM/service/message-service/internal/config"
	"github.com/go-redis/redis/v8"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Message{})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.CacheRedis.Host,
		Password: c.CacheRedis.Pass,
	})

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  redisClient,
	}
}
