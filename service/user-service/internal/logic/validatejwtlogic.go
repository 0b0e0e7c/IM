package logic

import (
	"context"
	"strconv"

	"github.com/0b0e0e7c/chat/component/auth"
	"github.com/0b0e0e7c/chat/service/user-service/internal/svc"
	"github.com/0b0e0e7c/chat/service/user-service/pb/user"
	"github.com/pkg/errors"

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
	// check if the token is in redis
	uidStr, err := l.svcCtx.Redis.Get(l.ctx, in.Token).Result()
	if err == nil {
		uid, err := strconv.ParseInt(uidStr, 10, 64)
		if err == nil {
			logx.Info("get uid with token from redis:", uid)
			return &user.ValidateResponse{
				Valid:  true,
				UserId: uid,
			}, nil
		}
		logx.Info("parse uid from redis failed")
	} else {
		logx.Info("token not found in redis")
	}

	valid, uid, err := auth.ValidateToken(in.Token)
	if errors.Is(err, auth.ErrTokenExpired) {
		logx.Info("token expired")
	}
	if err != nil {
		return nil, err
	}

	return &user.ValidateResponse{
			Valid:  valid,
			UserId: uid},
		nil
}
