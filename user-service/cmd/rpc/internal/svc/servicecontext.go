package svc

import (
	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/internal/config"
	"github.com/0b0e0e7c/IM/user-service/model"

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.CacheRedis.Host,
		Password: c.CacheRedis.Pass,
	})
	db.AutoMigrate(&model.User{})

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  redisClient,
	}
}
