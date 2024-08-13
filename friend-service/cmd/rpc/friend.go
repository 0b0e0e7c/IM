package main

import (
	"flag"
	"fmt"

	"github.com/0b0e0e7c/IM/friend-service/cmd/rpc/internal/config"
	"github.com/0b0e0e7c/IM/friend-service/cmd/rpc/internal/server"
	"github.com/0b0e0e7c/IM/friend-service/cmd/rpc/internal/svc"
	"github.com/0b0e0e7c/IM/friend-service/cmd/rpc/pb/friend"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "friend-service/cmd/rpc/etc/friend.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		friend.RegisterFriendServer(grpcServer, server.NewFriendServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
