package logic

import (
	"context"

	"IM/user-service/cmd/api/internal/svc"
	"IM/user-service/cmd/api/internal/types"
	"IM/user-service/model"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.UserRequest) (resp *types.UserResponse, err error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username or password is empty")
	}

	user := model.User{
		Username: req.Username,
		Password: hashing(req.Password),
	}

	result := l.svcCtx.DB.Create(&user)
	if result.Error != nil {
		logx.Error("create error:", result.Error)
		return nil, result.Error
	}

	resp = &types.UserResponse{
		Id:       int64(user.ID),
		Username: user.Username,
	}

	return resp, nil

}
