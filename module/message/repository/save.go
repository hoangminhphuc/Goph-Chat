package repository

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/module/message/model"
)

func (s *sqlRepo) SaveMessage(ctx context.Context, message *model.Message) error {
	db := s.db.Table(message.TableName())
	
	if err := db.Create(message).Error; err != nil {
		return err
	}
	return nil
}