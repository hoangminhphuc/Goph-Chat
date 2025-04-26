package repository

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/module/room/model"
)

func (s *sqlRepo) DeleteRoomByID(ctx context.Context, cond map[string]interface{}) error {
	db := s.db.Table(model.Room{}.TableName())

	if err := db.Where(cond).Delete(&model.Room{}).Error; err != nil {
		return common.ErrCannotDelete("room", err)
	}

	return nil
}