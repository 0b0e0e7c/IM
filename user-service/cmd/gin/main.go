package main

import (
	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/pb/user"

	"github.com/gin-gonic/gin"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

var rpcClient user.UserClient

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
	rpcClient = user.NewUserClient(EtcdClient.Conn())

	r := gin.Default()

	r.POST("/api/user/register", func(ctx *gin.Context) {
		Register(ctx, rpcClient)
	})
	r.POST("/api/user/login", func(ctx *gin.Context) {
		Login(ctx, rpcClient)
	})

	r.Run(":8888")
}
