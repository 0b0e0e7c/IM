package logic

import (
	"context"

	"github.com/0b0e0e7c/IM/service/message-service/internal/svc"
	"github.com/0b0e0e7c/IM/service/message-service/pb/msgservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessagesLogic) GetMessages(in *msgservice.GetMessagesRequest) (*msgservice.GetMessagesResponse, error) {
	// todo: add your logic here and delete this line

	return &msgservice.GetMessagesResponse{}, nil
}
