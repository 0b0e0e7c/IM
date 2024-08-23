package svc

import (
	"sync"

	"github.com/0b0e0e7c/chat/dao"
	"github.com/0b0e0e7c/chat/model"
	"github.com/0b0e0e7c/chat/service/user-service/internal/config"
	"github.com/go-redis/redis/v8"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client

	userDAO       *dao.UserDAO
	userDAOLoaded sync.Once
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.User{})

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

func (s *ServiceContext) GetUserDAO() *dao.UserDAO {
	s.userDAOLoaded.Do(func() {
		s.userDAO = dao.NewUserDAO(s.DB)
	})
	return s.userDAO
}
