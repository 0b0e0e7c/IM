package model

import (
	"time"
)

/* proto
message Message {
    int64 msg_id = 1;
    int64 sender_id = 2;
    int64 receiver_id = 3;
    string content = 4;
    int64 timestamp = 5;
}
*/

// Message 表结构
type Message struct {
	MsgID      int64     `gorm:"primaryKey;autoIncrement" json:"msg_id"`
	SenderId   int64     `gorm:"index;not null" json:"sender_id"`
	ReceiverId int64     `gorm:"index;not null" json:"receiver_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Timestamp  time.Time `gorm:"not null" json:"timestamp"`
}
