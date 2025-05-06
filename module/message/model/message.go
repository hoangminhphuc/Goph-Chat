package model

import (
	"github.com/hoangminhphuc/goph-chat/common/models"
	roomModel "github.com/hoangminhphuc/goph-chat/module/room/model"
	userModel "github.com/hoangminhphuc/goph-chat/module/user/model"
)

type Message struct {
	models.BaseModel
	Content string `json:"content" gorm:"column:content;"`
	RoomID  int    `json:"room_id" gorm:"column:room_id;"`
	UserID 	int    `json:"user_id" gorm:"column:user_id;"`
	Room 	*roomModel.Room  `json:"-" gorm:"foreignKey:room_id"`
	User 	*userModel.User `json:"-" gorm:"foreignKey:user_id"`
}

func (Message) TableName() string { return "messages" }