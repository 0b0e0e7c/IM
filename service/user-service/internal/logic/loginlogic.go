package logic

import (
	"context"
	"errors"
	"time"

	"github.com/0b0e0e7c/chat/component/auth"
	"github.com/0b0e0e7c/chat/model"
	"github.com/0b0e0e7c/chat/service/user-service/internal/svc"
	"github.com/0b0e0e7c/chat/service/user-service/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(req *user.LoginRequest) (*user.LoginResponse, error) {
	var loginUser model.User
	result := l.svcCtx.DB.Where("username = ? AND password = ?", req.Username, hashing(req.Username, req.Password)).First(&loginUser)
	if result.Error != nil {
		return nil, errors.New("invalid username or password")
	}

	// 生成JWT令牌
	token, err := auth.GenerateToken(int64(loginUser.ID), loginUser.Username)
	if err != nil {
		return nil, err
	}

	// 将JWT令牌存储到Redis中，并设置过期时间
	err = l.svcCtx.Redis.Set(l.ctx, token, int64(loginUser.ID), 24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	resp := &user.LoginResponse{
		UserId:   int64(loginUser.ID),
		Username: loginUser.Username,
		Token:    token,
		Success:  true,
	}

	return resp, nil
}
