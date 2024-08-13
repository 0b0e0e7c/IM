package logic

import (
	"context"

	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/internal/svc"
	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/pb/user"
	"github.com/0b0e0e7c/IM/user-service/component/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateJWTLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateJWTLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateJWTLogic {
	return &ValidateJWTLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ValidateJWTLogic) ValidateJWT(in *user.ValidateRequest) (resp *user.ValidateResponse, err error) {
	valid, err := auth.ValidateToken(in.Token)

	if err != nil {
		return nil, err
	}

	resp = &user.ValidateResponse{
		Valid: valid,
	}
	return
}
