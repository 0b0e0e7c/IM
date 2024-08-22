package logic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/0b0e0e7c/chat/model"
	"github.com/0b0e0e7c/chat/service/user-service/internal/svc"
	"github.com/0b0e0e7c/chat/service/user-service/pb/user"

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

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	if in.Username == "" || in.Password == "" {
		return nil, errors.New("username or password is empty")
	}

	newUser := model.User{
		Username: in.Username,
		Password: hashing(in.Username, in.Password),
	}

	result := l.svcCtx.DB.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	resp := &user.RegisterResponse{
		UserId:   int64(newUser.ID),
		Username: newUser.Username,
		Success:  true,
	}

	return resp, nil
}

func hashing(username, password string) string {
	hash := sha256.New()

	hash.Write([]byte(username + password))

	return hex.EncodeToString(hash.Sum(nil))
}
