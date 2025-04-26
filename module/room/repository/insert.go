package repository

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/module/room/dto"
)

func (s *sqlRepo) CreateRoom(ctx context.Context, data *dto.RoomCreation) error {
	db := s.db.Table(data.TableName())

	if err := db.Create(data).Error; err != nil {
		return err
	}
	
	return nil
}