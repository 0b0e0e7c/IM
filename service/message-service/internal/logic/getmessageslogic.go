package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0b0e0e7c/IM/component/common"
	"github.com/0b0e0e7c/IM/model"
	"github.com/0b0e0e7c/IM/service/message-service/internal/svc"
	"github.com/0b0e0e7c/IM/service/message-service/pb/message"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessagesLogic) GetMessages(in *message.GetMessagesRequest) (*message.GetMessagesResponse, error) {
	lowerID, higherID := common.LowHigh(in.UserId, in.PeerId)

	if in.Limit == 0 {
		in.Limit = 100
	}

	// 从 Redis 获取消息
	messages, err := l.pullMsgFromRedis(lowerID, higherID, in.Offset, in.Offset+in.Limit-1)
	if err != nil {
		return nil, err
	}

	// 从 MySQL 中获取
	if len(messages) == 0 || len(messages) < int(in.Limit) {
		messages, err = l.pullMsgFromMySQL(lowerID, higherID, in.Offset, in.Limit)
		if err != nil {
			return nil, err
		}
		// 将从 MySQL 获取的消息缓存到 Redis
		l.pushMsgToRedis(lowerID, higherID, messages)
	}

	var msgList []*message.Message
	for _, msg := range messages {
		msgList = append(msgList, &message.Message{
			MsgId:      msg.MsgID,
			SenderId:   msg.SenderId,
			ReceiverId: msg.ReceiverId,
			Content:    msg.Content,
			Timestamp:  msg.Timestamp.Unix(),
		})
	}

	return &message.GetMessagesResponse{
		Messages: msgList,
	}, nil
}

func (l *GetMessagesLogic) pullMsgFromMySQL(senderId, receiverId, offset, limit int64) (messages []*model.Message, err error) {
	logx.Infof("pull messages from MySQL")
	// 查询 MySQL 数据库中的消息记录
	err = l.svcCtx.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderId, receiverId, receiverId, senderId).
		Order("timestamp DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return
}

func (l *GetMessagesLogic) pullMsgFromRedis(lowerID, higherID, start, end int64) (messages []*model.Message, err error) {
	logx.Infof("pull messages from Redis")
	key := fmt.Sprintf("chat:%d-%d", lowerID, higherID)

	// 使用 ZREVRANGE 获取有序集合中的消息
	msgs, err := l.svcCtx.Redis.ZRevRange(l.ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}

	// 遍历每个消息字符串，并反序列化为 Message 结构体
	for _, msgStr := range msgs {
		var msg *model.Message
		if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	logx.Infof("messages from Redis: %+v len: %d", messages, len(messages))

	return messages, nil
}

func (l *GetMessagesLogic) pushMsgToRedis(lowerID, higherID int64, messages []*model.Message) error {

	logx.Infof("push missing message to Redis: %+v len: %d", messages, len(messages))

	key := fmt.Sprintf("chat:%d-%d", lowerID, higherID)

	for _, msg := range messages {
		msgEntry, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		// 使用 ZADD 将消息添加到有序集合中
		err = l.svcCtx.Redis.ZAdd(l.ctx, key, &redis.Z{
			Score:  float64(msg.MsgID),
			Member: string(msgEntry),
		}).Err()
		if err != nil {
			return err
		}

	}

	// 保留最新的10条消息
	err := l.svcCtx.Redis.ZRemRangeByRank(l.ctx, key, 0, -11).Err()
	if err != nil {
		return err
	}

	return nil
}
