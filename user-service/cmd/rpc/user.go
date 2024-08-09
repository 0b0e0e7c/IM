package main

import (
	"flag"
	"fmt"

	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/internal/config"
	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/internal/server"
	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/internal/svc"
	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "user-service/cmd/rpc/etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
