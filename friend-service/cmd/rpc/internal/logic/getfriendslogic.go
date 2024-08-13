package logic

import (
	"context"

	"github.com/0b0e0e7c/IM/friend-service/cmd/rpc/internal/svc"
	"github.com/0b0e0e7c/IM/friend-service/cmd/rpc/pb/friend"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsLogic {
	return &GetFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendsLogic) GetFriends(in *friend.GetFriendsRequest) (*friend.GetFriendsResponse, error) {
	// todo: add your logic here and delete this line

	return &friend.GetFriendsResponse{}, nil
}
