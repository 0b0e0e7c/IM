package logic

import (
	"context"
	"errors"

	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/internal/svc"
	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/pb/user"
	"github.com/0b0e0e7c/IM/user-service/model"

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

func (l *LoginLogic) Login(req *user.UserRequest) (resp *user.UserResponse, err error) {
	var loginUser model.User
	result := l.svcCtx.DB.Where("username = ? AND password = ?", req.Username, hashing(req.Password)).First(&loginUser)
	if result.Error != nil {
		return nil, errors.New("invalid username or password")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	resp = &user.UserResponse{
		Id:       int64(loginUser.ID),
		Username: loginUser.Username,
	}

	return resp, nil
}
