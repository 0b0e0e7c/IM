package logic

import (
	"context"
	"fmt"

	"github.com/0b0e0e7c/IM/friend-service/cmd/rpc/internal/svc"
	"github.com/0b0e0e7c/IM/friend-service/cmd/rpc/pb/friend"
	"github.com/0b0e0e7c/IM/friend-service/model"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendError struct {
	Message string
}

func (e *AddFriendError) Error() string {
	return fmt.Sprintf("add friend error: %s", e.Message)
}

type AddFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddFriendLogic) AddFriend(in *friend.AddFriendRequest) (resp *friend.AddFriendResponse, err error) {
	if in.UserId == in.FriendId {
		return nil, &AddFriendError{Message: "can't add yourself as friend"}
	}

	if in.UserId > in.FriendId {
		in.UserId, in.FriendId = in.FriendId, in.UserId
	}

	// check if user and friend exists
	var existingUser1, existingUser2 model.User
	if err := l.svcCtx.DB.Where("id = ?", in.UserId).First(&existingUser1).Error; err != nil {
		return nil, &AddFriendError{Message: "user not found"}
	}
	if err := l.svcCtx.DB.Where("id = ?", in.FriendId).First(&existingUser2).Error; err != nil {
		return nil, &AddFriendError{Message: "friend not found"}
	}

	var existingFriend model.Friend
	resp = &friend.AddFriendResponse{}
	if err := l.svcCtx.DB.Where("user_id = ? AND friend_id = ?", in.UserId, in.FriendId).First(&existingFriend).Error; err == gorm.ErrRecordNotFound {
		friend := model.Friend{
			UserID:   in.UserId,
			FriendID: in.FriendId,
			Status:   1, // 默认状态为已添加
		}
		result := l.svcCtx.DB.Create(&friend)
		if result.Error != nil {
			return nil, result.Error
		} else {
			resp.Success = true
		}
	} else {
		if existingFriend.Status == 1 {
			return nil, &AddFriendError{Message: "friend already added"}
		}

		result := l.svcCtx.DB.Model(&model.Friend{}).Where("id = ?", existingFriend.ID).Update("status", 1)
		if result.Error != nil {
			return nil, result.Error
		} else {
			resp.Success = false
		}
	}

	return resp, nil
}
