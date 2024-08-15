package logic

import (
	"context"

	"github.com/0b0e0e7c/IM/model"
	"github.com/0b0e0e7c/IM/service/friend-service/internal/svc"
	"github.com/0b0e0e7c/IM/service/friend-service/pb/friend"

	"github.com/zeromicro/go-zero/core/logx"
)

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

func (l *AddFriendLogic) AddFriend(in *friend.AddFriendRequest) (*friend.AddFriendResponse, error) {
	if in.UserId == in.FriendId {
		return nil, &FriendServiceError{Message: "can't add yourself as friend"}
	}

	if in.UserId > in.FriendId {
		in.UserId, in.FriendId = in.FriendId, in.UserId
	}

	// check if user and friend exists
	var existingUsers []model.User
	if err := l.svcCtx.DB.Where("id IN (?)", []int64{in.UserId, in.FriendId}).Find(&existingUsers).Error; err != nil || len(existingUsers) != 2 {
		return nil, &FriendServiceError{Message: "user or friend not found"}
	}

	// 查询是否已有好友关系
	var existingFriend model.Friend
	if err := l.svcCtx.DB.Where("user_id = ? AND friend_id = ?", in.UserId, in.FriendId).First(&existingFriend).Error; err == nil {
		// 已存在好友关系
		if existingFriend.Status == 1 {
			return &friend.AddFriendResponse{Success: false}, &FriendServiceError{Message: "friend already added"}
		}
		// 未来可能有删除好友的情况，此时更新状态为已添加
		if err := l.svcCtx.DB.Model(&existingFriend).Update("status", 1).Error; err != nil {
			return nil, err
		}

		return &friend.AddFriendResponse{Success: true}, nil
	}

	// 如果不存在，则创建新的好友关系
	newFriend := model.Friend{
		UserID:   in.UserId,
		FriendID: in.FriendId,
		Status:   1, // 设置1为已添加
	}
	if err := l.svcCtx.DB.Create(&newFriend).Error; err != nil {
		return nil, err
	}

	return &friend.AddFriendResponse{Success: true}, nil
}
