package repository

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/module/user/model"
	"gorm.io/gorm"
)

func (s *sqlRepo) FindUser(ctx context.Context, conditions map[string]interface{}) (*model.User, error) {
	db := s.db.Table(model.User{}.TableName())
	var user model.User
	
	if err := db.Where(conditions).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &user, nil
}