package model

import "github.com/hoangminhphuc/goph-chat/common/models"

type User struct {
	models.BaseModel
	Email    string   `json:"email" gorm:"column:email;"`
	Password string   `json:"password" gorm:"column:password;"`
	LastName string   `json:"last_name" gorm:"column:last_name;"`
	FirstName string  `json:"first_name" gorm:"column:first_name;"`
}