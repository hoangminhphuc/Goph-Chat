package dto

import "github.com/hoangminhphuc/goph-chat/module/room/model"

type RoomCreation struct {
	ID          int    `json:"id" gorm:"column:id;"`
	Name        string `json:"name" gorm:"column:name;"`
	Description string `json:"description" gorm:"column:description;"`
	OwnerID     int    `json:"user_id" gorm:"column:user_id;"`
}

func (RoomCreation) TableName() string { return model.Room{}.TableName() }

func (room *RoomCreation) Mask() {
	room.Mask()
}
