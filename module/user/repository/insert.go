package repository

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/module/user/dto"
)

func (s *sqlRepo) CreateUser(ctx context.Context, data *dto.UserRegister) error {
	db := s.db.Begin()
	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback()
		return err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return err
	}
	return nil

}