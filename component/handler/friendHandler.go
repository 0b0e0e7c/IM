package handler

import (
	"context"
	"net/http"

	"github.com/0b0e0e7c/IM/service/friend-service/pb/friend"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

func AddFriend(c *gin.Context, client friend.FriendServiceClient) {
	var req struct {
		FriendID int64 `json:"friend_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	logx.Info("userID: ", userID)

	_, err := client.AddFriend(context.Background(), &friend.AddFriendRequest{
		UserId:   userID.(int64),
		FriendId: req.FriendID,
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
		"status": "success",
	})
}

func GetFriends(c *gin.Context, client friend.FriendServiceClient) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	resp, err := client.GetFriends(context.Background(), &friend.GetFriendsRequest{
		UserId: userID.(int64),
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
		"friends": resp.FriendIds,
		"status":  "success",
	})
}
