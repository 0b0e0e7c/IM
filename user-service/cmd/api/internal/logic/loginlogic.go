package logic

import (
	"context"

	"IM/user-service/cmd/api/internal/svc"
	"IM/user-service/cmd/api/internal/types"
	"IM/user-service/model"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserRequest) (resp *types.UserResponse, err error) {
	var user model.User
	result := l.svcCtx.DB.Where("username = ? AND password = ?", req.Username, hashing(req.Password)).First(&user)
	if result.Error != nil {
		return nil, errors.New("invalid username or password")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	resp = &types.UserResponse{
		Id:       int64(user.ID),
		Username: user.Username,
	}

	return resp, nil
}
