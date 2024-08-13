package model

import (
	"gorm.io/gorm"
)

// Friend 表结构
type Friend struct {
	gorm.Model
	UserID   int64 `gorm:"not null" json:"user_id"`
	FriendID int64 `gorm:"not null" json:"friend_id"`
	Status   int32 `gorm:"not null" json:"status"`
}
