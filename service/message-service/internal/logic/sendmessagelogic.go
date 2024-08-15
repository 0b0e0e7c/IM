package logic

import (
	"context"

	"github.com/0b0e0e7c/IM/service/message-service/internal/svc"
	"github.com/0b0e0e7c/IM/service/message-service/pb/msgservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageLogic) SendMessage(in *msgservice.SendMessageRequest) (*msgservice.SendMessageResponse, error) {
	// todo: add your logic here and delete this line

	return &msgservice.SendMessageResponse{}, nil
}
