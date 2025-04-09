package model

import "github.com/hoangminhphuc/goph-chat/common/models"

type Room struct {
	models.BaseModel
	Name 				string `json:"name" gorm:"column:name;"`
	Description string `json:"description" gorm:"column:description;"`
}

func (Room) TableName() string { return "rooms" }
