package repository

import (
	"context"
	"net/http"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/message/model"
)

func (s *sqlRepo) ListMessage(ctx context.Context, roomID int, paging *utils.Paging) ([]model.Message, error) {
	db := s.db.Table(model.Message{}.TableName()).Where("room_id = ?", roomID)

	if err := db.Select("id").Count(&paging.Total).Error; err != nil {
		return nil, common.WrapError(err, "cannot count total messages", http.StatusBadRequest)
	}

	if paging.FakeCursor != "" {
		uuid, err := utils.DecodeID(paging.FakeCursor)
		if err != nil {
			return nil, common.WrapError(err, "invalid cursor", http.StatusBadRequest)
		}
		db = db.Where("id < ?", uuid)
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	var messages []model.Message

	if err := db.Select("*").Order("id desc").Limit(paging.Limit).Find(&messages).Error; err != nil {
		return nil, common.WrapError(err, "cannot find messages", http.StatusBadRequest)
	}

	if len(messages) > 0 {
		messages[len(messages)-1].Mask()
		paging.NextCursor = messages[len(messages)-1].FakeID
	}

	return messages, nil
}
