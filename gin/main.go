package main

import (
	"github.com/0b0e0e7c/IM/service/user-service/pb/user"
	"github.com/gin-gonic/gin"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

func main() {
	// 创建 RPC 客户端配置
	EtcdClient, err := zrpc.NewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "user.rpc",
		},
	})
	if err != nil {
		panic(err)
	}
	userRPCClient := user.NewUserServiceClient(EtcdClient.Conn())

	r := gin.Default()

	r.POST("/api/user/register", func(ctx *gin.Context) {
		Register(ctx, userRPCClient)
	})
	r.POST("/api/user/login", func(ctx *gin.Context) {
		Login(ctx, userRPCClient)
	})
	r.POST("/api/user/ValidateJWT", func(ctx *gin.Context) {
		ValidateJWT(ctx, userRPCClient)
	})

	r.Run(":8888")
}
