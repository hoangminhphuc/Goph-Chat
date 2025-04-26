package repository

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/module/room/model"
	"gorm.io/gorm"
)

func (s *sqlRepo) FindRoom(ctx context.Context, cond map[string]interface{}) (*model.Room, error) {
	db := s.db.Table(model.Room{}.TableName())

	var roomData model.Room

	if err := db.Where(cond).First(&roomData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &roomData, nil
}