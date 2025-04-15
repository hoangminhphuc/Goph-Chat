package dto

import "github.com/hoangminhphuc/goph-chat/module/user/model"

type UserRegister struct {
	Email 			string 		`json:"email" gorm:"column:email;"`
	Password 		string 		`json:"password" gorm:"column:password;"`
	Salt 				string   	`json:"-" gorm:"column:salt;"`
	LastName 		string   	`json:"last_name" gorm:"column:last_name;"`
	FirstName 	string  	`json:"first_name" gorm:"column:first_name;"`
	Phone    		string   	`json:"phone" gorm:"column:phone;"` 
	Role 				string 	`json:"-" gorm:"column:role;"`
}

func (UserRegister) TableName() string { return model.User{}.TableName() } 

type UserLogin struct {
	Email string `json:"email" gorm:"column:email;"`
	Password string `json:"password" gorm:"column:password;"`
}

func (UserLogin) TableName() string { return model.User{}.TableName() }