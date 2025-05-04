package repository

import (
	"context"
	"fmt"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/internal/cache"
	"github.com/hoangminhphuc/goph-chat/module/user/model"
	"gorm.io/gorm"
)

func (s *sqlRepo) FindUser(ctx context.Context, conditions map[string]interface{}) (*model.User, error) {
	db := s.db.Table(model.User{}.TableName())

	var key string
	if id, ok := conditions["id"].(int); ok {
		key = fmt.Sprintf("user:id:%d", id)
	} else if email, ok := conditions["email"].(string); ok {
		key = fmt.Sprintf("user:email:%s", email)
	}
	
	var user model.User

	cache := cache.GetCache()
	if err := cache.Get(ctx, key, &user); err == nil {
		return &user, nil
	}
	
	if err := db.Where(conditions).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	_ = cache.Set(ctx, key, &user, common.UserProfileTTL)

	return &user, nil
}