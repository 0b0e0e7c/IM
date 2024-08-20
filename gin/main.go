package main

import (
	"github.com/0b0e0e7c/IM/component/handler"
	"github.com/0b0e0e7c/IM/component/middleware"
	"github.com/0b0e0e7c/IM/service/friend-service/pb/friend"
	"github.com/0b0e0e7c/IM/service/user-service/pb/user"

	"github.com/gin-gonic/gin"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

func main() {
	// 创建 RPC 客户端配置
	userClient, err := zrpc.NewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "user.rpc",
		},
	})
	if err != nil {
		panic(err)
	}
	userRPCClient := user.NewUserServiceClient(userClient.Conn())

	r := gin.Default()

	r.POST("/api/user/register", func(ctx *gin.Context) {
		handler.Register(ctx, userRPCClient)
	})
	r.POST("/api/user/login", func(ctx *gin.Context) {
		handler.Login(ctx, userRPCClient)
	})

	r.POST("/api/user/ValidateJWT", func(ctx *gin.Context) {
		handler.ValidateJWT(ctx, userRPCClient)
	})

	friendClient, err := zrpc.NewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "friend.rpc",
		},
	})
	if err != nil {
		panic(err)
	}
	fiendRPCClient := friend.NewFriendServiceClient(friendClient.Conn())

	authGroup := r.Group("/api/friend")
	authGroup.Use(middleware.JWTMiddleware(userRPCClient))
	{
		authGroup.POST("/add", func(ctx *gin.Context) {
			handler.AddFriend(ctx, fiendRPCClient)
		})
		authGroup.POST("/get", func(ctx *gin.Context) {
			handler.GetFriends(ctx, fiendRPCClient)
		})
	}

	r.Run(":8888")
}
