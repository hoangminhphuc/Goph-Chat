package module

import (
	// userModel "github.com/hoangminhphuc/goph-chat/module/user/model"
	roomModel "github.com/hoangminhphuc/goph-chat/module/room/model"
)

func GetAllModels() []interface{} {
	return []interface{}{ 
		// &userModel.User{},
		&roomModel.Room{},
	}
}