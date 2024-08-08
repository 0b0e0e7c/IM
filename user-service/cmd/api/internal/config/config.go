package config

import (
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type Config struct {
	rest.RestConf
	DataSource string
}

var DB *gorm.DB
