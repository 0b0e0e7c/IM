package logic

import (
	"context"

	"github.com/0b0e0e7c/IM/component/auth"
	"github.com/0b0e0e7c/IM/service/user-service/internal/svc"
	"github.com/0b0e0e7c/IM/service/user-service/pb/user"

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

func (l *ValidateJWTLogic) ValidateJWT(in *user.ValidateRequest) (*user.ValidateResponse, error) {
	valid, uid, err := auth.ValidateToken(in.Token)

	if err != nil {
		return nil, err
	}

	return &user.ValidateResponse{
			Valid:  valid,
			UserId: uid},
		nil
}
