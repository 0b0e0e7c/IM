package logic

import (
	"context"
	"time"

	"github.com/0b0e0e7c/chat/component/auth"
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
	hashedPassword := Hashing(req.Username, req.Password)

	userDAO := l.svcCtx.GetUserDAO()

	loginUser, err := userDAO.FindUserByUsernameAndPassword(req.Username, hashedPassword)
	if err != nil {
		return nil, err
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

	// 返回登录响应
	resp := &user.LoginResponse{
		UserId:   int64(loginUser.ID),
		Username: loginUser.Username,
		Token:    token,
		Success:  true,
	}

	return resp, nil
}
