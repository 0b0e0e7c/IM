package handler

import (
	"context"
	"net/http"

	"github.com/0b0e0e7c/IM/service/message-service/pb/message"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

func SendMsg(c *gin.Context, client message.MessageServiceClient) {
	var req struct {
		ReceiverID int64  `json:"receiver_id" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	logx.Info("userID: ", userID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	resp, err := client.SendMessage(context.Background(), &message.SendMessageRequest{
		SenderId:   userID.(int64),
		ReceiverId: req.ReceiverID,
		Content:    req.Content,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"msg":    err.Error(),
			})
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"msg":    st.Message(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": resp.Success,
		"msg":    "success",
	})
}

type Message struct {
	MsgId      int64  `json:"msg_id"`
	SenderId   int64  `json:"sender_id"`
	ReceiverId int64  `json:"receiver_id"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
}

func GetMsg(c *gin.Context, client message.MessageServiceClient) {
	var req struct {
		PeerId int64 `json:"peer_id" binding:"required"`
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	logx.Info("userID: ", userID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	resp, err := client.GetMessages(context.Background(), &message.GetMessagesRequest{
		UserId: userID.(int64),
		PeerId: req.PeerId,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"msg":    err.Error(),
			})
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"msg":    st.Message(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  "success",
		"messages": resp.Messages,
	})
}
