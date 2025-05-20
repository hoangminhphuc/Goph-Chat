package dto

import "github.com/hoangminhphuc/goph-chat/module/message/model"

type RoomInfo struct {
	RoomID	 int 	`form:"roomID"`  
}

type MessageUpdate struct {
	RoomID	 int    
	Content string `json:"content" gorm:"column:content;"`
	UserID	 int   
}

func (MessageUpdate) TableName() string {
	return model.Message{}.TableName()
}