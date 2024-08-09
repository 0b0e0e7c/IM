package logic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/internal/svc"
	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/pb/user"
	"github.com/0b0e0e7c/IM/user-service/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.UserRequest) (resp *user.UserResponse, err error) {
	if in.Username == "" || in.Password == "" {
		return nil, errors.New("username or password is empty")
	}

	newUser := model.User{
		Username: in.Username,
		Password: hashing(in.Password),
	}

	result := l.svcCtx.DB.Create(&newUser)
	if result.Error != nil {
		logx.Error("create error:", result.Error)
		return nil, result.Error
	}

	resp = &user.UserResponse{
		Id:       int64(newUser.ID),
		Username: newUser.Username,
	}

	return resp, nil
}

func hashing(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
