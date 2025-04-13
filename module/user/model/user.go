package model

import "github.com/hoangminhphuc/goph-chat/common/models"

type UserRole int

const (
	RoleUser UserRole = iota
	RoleAdmin 

)

// Fixed length
func (r UserRole) String() string {
	return [...]string{"user", "admin"}[r]
}

type User struct {
	models.BaseModel
	Email    		string   	`json:"email" gorm:"column:email;"`
	Password 		string   	`json:"-" gorm:"column:password;"`
	Salt 				string   	`json:"-" gorm:"column:salt;"`
	LastName 		string   	`json:"last_name" gorm:"column:last_name;"`
	FirstName 	string  	`json:"first_name" gorm:"column:first_name;"`
	Phone    		string   	`json:"phone" gorm:"column:phone;"`
	Role 				UserRole 	`json:"role" gorm:"column:role;"`
	Status  		string   	`json:"status" gorm:"column:status;"`
}
