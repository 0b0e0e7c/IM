package logic

import (
	"context"
	"errors"

	"github.com/0b0e0e7c/chat/model"
	"github.com/0b0e0e7c/chat/service/friend-service/internal/svc"
	"github.com/0b0e0e7c/chat/service/friend-service/pb/friend"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsLogic {
	return &GetFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendsLogic) GetFriends(in *friend.GetFriendsRequest) (*friend.GetFriendsResponse, error) {
	user := model.User{}

	// 检查用户是否存在
	if err := l.svcCtx.DB.Where("id = ?", in.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &FriendServiceError{Message: "user not found"}
		}
		return nil, err
	}

	// 查找 friend 表中 user_id 或 friend_id 为 in.UserId 且 status 为 1 的记录，返回所有 friend_id
	var friends []model.Friend
	if err := l.svcCtx.DB.Where("(user_id = ? OR friend_id = ?) AND status = 1", in.UserId, in.UserId).Find(&friends).Error; err != nil {
		return nil, err
	}

	// 构建响应
	resp := &friend.GetFriendsResponse{}
	for _, friend := range friends {
		if friend.UserID == in.UserId {
			resp.FriendIds = append(resp.FriendIds, friend.FriendID)
		} else {
			resp.FriendIds = append(resp.FriendIds, friend.UserID)
		}
	}

	return resp, nil
}
