package model

import (
	"fmt"
	"net/http"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/models"
)

type UserRole int

const (
	RoleUser UserRole = iota
	RoleAdmin 

)

// Fixed length
func (r UserRole) String() string {
	return [...]string{"user", "admin"}[r]
}

// Implements Scanner interface for GORM to use
func (role *UserRole) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
			return common.NewError(fmt.Sprintf("Failed to map value: %v", value),
				http.StatusInternalServerError)
		}

	var r UserRole
	roleValue := string(bytes)

	if roleValue == "user" {
			r = RoleUser
	} else if roleValue == "admin" {
			r = RoleAdmin
	}

	*role = r
	return nil
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

func (User) TableName() string { return "users" }
