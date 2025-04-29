package business

import (
	"context"
	"net/http"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/module/room/model"
)

type DeleteRoomRepo interface {
	FindRoom(ctx context.Context, cond map[string]interface{}) (*model.Room, error)
	DeleteRoomByID(ctx context.Context, cond map[string]interface{}) error
}

type DeleteRoomBusiness struct {
	repo DeleteRoomRepo
}

func NewDeleteRoomBusiness(repo DeleteRoomRepo) *DeleteRoomBusiness {
	return &DeleteRoomBusiness{repo: repo}
}

func (dr *DeleteRoomBusiness) DeleteRoomByID(ctx context.Context, id int) error {
	data, err := dr.repo.FindRoom(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return common.WrapError(err, "cannot find room", http.StatusInternalServerError)
	}

	if data.DeletedAt.Valid {
		return common.NewError("room already deleted", http.StatusBadRequest)
	}

	return dr.repo.DeleteRoomByID(ctx, map[string]interface{}{"id": id})
}