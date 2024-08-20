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
	userRPCClient := initUserRPCClient()
	friendRPCClient := initFriendRPCClient()

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

	authGroup := r.Group("/api/friend")
	authGroup.Use(middleware.JWTMiddleware(userRPCClient))
	{
		authGroup.POST("/add", func(ctx *gin.Context) {
			handler.AddFriend(ctx, friendRPCClient)
		})
		authGroup.POST("/get", func(ctx *gin.Context) {
			handler.GetFriends(ctx, friendRPCClient)
		})
	}

	r.Run(":8888")
}

func initUserRPCClient() user.UserServiceClient {
	userClient, err := zrpc.NewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "user.rpc",
		},
	})
	if err != nil {
		panic(err)
	}
	return user.NewUserServiceClient(userClient.Conn())
}

func initFriendRPCClient() friend.FriendServiceClient {
	friendClient, err := zrpc.NewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "friend.rpc",
		},
	})
	if err != nil {
		panic(err)
	}
	return friend.NewFriendServiceClient(friendClient.Conn())
}
