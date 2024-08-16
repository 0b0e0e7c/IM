package logic

import (
	"context"
	"strconv"
	"time"

	"github.com/0b0e0e7c/IM/model"
	"github.com/0b0e0e7c/IM/service/message-service/internal/svc"
	"github.com/0b0e0e7c/IM/service/message-service/pb/msgservice"

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

func (l *SendMessageLogic) SendMessage(in *msgservice.SendMessageRequest) (*msgservice.SendMessageResponse, error) {

	u1, u2 := compareAndSwap(in.SenderId, in.ReceiverId)

	if err := l.svcCtx.DB.Where("(user_id = ? AND friend_id = ? )  AND status = 1", u1, u2).First(&model.Friend{}).Error; err != nil {
		return nil, &MessageServiceError{Message: "sender and receiver are not friends"}
	}

	// create message
	message := model.Message{
		SenderId:   in.SenderId,
		ReceiverId: in.ReceiverId,
		Content:    in.Content,
		Timestamp:  time.Now(),
	}

	if err := l.svcCtx.DB.Create(&message).Error; err != nil {
		return nil, &MessageServiceError{Message: "failed to create message"}
	}

	l.pushMsgToRedis(message.MsgID, message.SenderId, message.ReceiverId, message.Content, message.Timestamp)

	return &msgservice.SendMessageResponse{
		Success: true,
	}, nil
}

func (l *SendMessageLogic) pushMsgToRedis(msgID, senderID, receiverID int64, content string, timestamp time.Time) error {
	// make sure the key is always the same
	var key string
	lowerID, higherID := compareAndSwap(senderID, receiverID)

	key = "chat:" + strconv.FormatInt(lowerID, 10) + ":" + strconv.FormatInt(higherID, 10)

	// msgID|senderID|receiverID|content|timestamp
	msgEntry := strconv.FormatInt(msgID, 10) + "|" +
		strconv.FormatInt(senderID, 10) + "|" +
		strconv.FormatInt(receiverID, 10) + "|" +
		content + "|" +
		timestamp.Format("2006-01-02 15:04:05")

	err := l.svcCtx.Redis.LPush(l.ctx, key, msgEntry).Err()
	if err != nil {
		return err
	}

	return nil
}

func compareAndSwap(a, b int64) (int64, int64) {
	if a > b {
		return b, a
	}
	return a, b
}
