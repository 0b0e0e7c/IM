package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"IM/user-service/cmd/api/internal/config"
	"IM/user-service/cmd/api/internal/handler"
	"IM/user-service/cmd/api/internal/svc"

	"github.com/gin-gonic/gin"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "user-service/cmd/api/etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// parse relative path to absolute path
	absConfigFile, err := filepath.Abs(*configFile)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}
	conf.MustLoad(absConfigFile, &c)

	ctx := svc.NewServiceContext(c)

	r := gin.Default()
	handler.RegisterHandlers(r, ctx)

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	http.ListenAndServe(addr, r)
}
