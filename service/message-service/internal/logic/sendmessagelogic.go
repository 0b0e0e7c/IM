package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0b0e0e7c/IM/component/common"
	"github.com/0b0e0e7c/IM/model"
	"github.com/0b0e0e7c/IM/service/message-service/internal/svc"
	"github.com/0b0e0e7c/IM/service/message-service/pb/message"
	"github.com/go-redis/redis/v8"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageLogic) SendMessage(in *message.SendMessageRequest) (*message.SendMessageResponse, error) {

	lowerID, higherID := common.LowHigh(in.SenderId, in.ReceiverId)

	if err := l.svcCtx.DB.Where("(user_id = ? AND friend_id = ? )  AND status = 1", lowerID, higherID).First(&model.Friend{}).Error; err != nil {
		return nil, &MessageServiceError{Message: "sender and receiver are not friends"}
	}

	// create message
	newMsg := model.Message{
		SenderId:   in.SenderId,
		ReceiverId: in.ReceiverId,
		Content:    in.Content,
		Timestamp:  time.Now(),
	}

	if err := l.svcCtx.DB.Create(&newMsg).Error; err != nil {
		return nil, &MessageServiceError{Message: "failed to create message"}
	}

	logx.Infof("new message created: %+v", newMsg)

	l.pushMsgToRedis(&newMsg)

	return &message.SendMessageResponse{
		Success: true,
	}, nil
}

func (l *SendMessageLogic) pushMsgToRedis(newMsg *model.Message) error {
	lowerID, higherID := common.LowHigh(newMsg.SenderId, newMsg.ReceiverId)

	key := fmt.Sprintf("chat:%d-%d", lowerID, higherID)

	msgEntry, err := json.Marshal(newMsg)
	if err != nil {
		return err
	}

	logx.Infof("put new message to Redis: %+v", newMsg)

	// 使用 ZADD 将消息添加到有序集合中
	err = l.svcCtx.Redis.ZAdd(l.ctx, key, &redis.Z{
		Score:  float64(newMsg.MsgID),
		Member: string(msgEntry),
	}).Err()
	if err != nil {
		return err
	}

	// 保留最新的10条消息
	err = l.svcCtx.Redis.ZRemRangeByRank(l.ctx, key, 0, -11).Err()
	if err != nil {
		return err
	}

	return nil
}
