package model

import (
	"github.com/hoangminhphuc/goph-chat/common/models"
	"github.com/hoangminhphuc/goph-chat/module/user/model"
	"gorm.io/gorm"
)

type Room struct {
	models.BaseModel
	Name 				string 							`json:"name" gorm:"column:name;"`
	Description string 							`json:"description" gorm:"column:description;"`
	OwnerID 		int 								`json:"-" gorm:"column:user_id;"`
	User 				*model.User 				`json:"owner" gorm:"foreignKey:OwnerID"` // Has Many relation
	DeletedAt 	gorm.DeletedAt 			`json:"-" gorm:"index"`
}

func (Room) TableName() string { return "rooms" }

func (room *Room) Mask() {
	room.BaseModel.Mask()

	if value := room.User; value != nil {
		value.Mask()
	}
}
