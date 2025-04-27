package repository

import (
	"context"
	"net/http"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/room/model"
)

func (s *sqlRepo) ListRoom(ctx context.Context, paging *utils.Paging) ([]model.Room, error) {
	db := s.db.Table(model.Room{}.TableName())

	/*
		We Count to know the total number of matching items,
		even if we are only fetching a few items now because of pagination.

		For example, the user may only see 5 items now (limit = 5). 
		But maybe there are 200 items total matching the search.

		Frontend can show: ➔ "Showing 5 of 200 items"
	*/

	/* 
		When counting, we don't need to load all columns from the table, just id col.
	*/
	if err := db.Select("id").Count(&paging.Total).Error; err != nil {
		return nil, common.WrapError(err, "cannot count total rooms", http.StatusBadRequest)
	}

	db = db.Scopes(PreloadScope("User"))

	if paging.FakeCursor != "" {
		uuid, err := utils.DecodeID(paging.FakeCursor)
		if err != nil {
			return nil, common.WrapError(err, "invalid cursor", http.StatusBadRequest)
		}

		db = db.Where("id < ?", uuid)
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	var rooms []model.Room

	if err :=db.Select("*").Order("id desc").
		Limit(paging.Limit).Find(&rooms).Error; err != nil {
			return nil, common.WrapError(err, "cannot find room", http.StatusBadRequest)
		}

	if len(rooms) > 0 {
		rooms[len(rooms) - 1].Mask()
		paging.NextCursor = rooms[len(rooms) - 1].FakeID
	}

	return rooms, nil
}