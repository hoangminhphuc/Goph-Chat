package business

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/module/room/model"
)

type GetRoomRepo interface {
	FindRoom(ctx context.Context, cond map[string]interface{}) (*model.Room, error)
}

type GetRoomBusiness struct {
	repo GetRoomRepo
}

func NewGetRoomBusiness(repo GetRoomRepo) *GetRoomBusiness {
	return &GetRoomBusiness{repo: repo}
}

func (gr *GetRoomBusiness) GetRoom(ctx context.Context, id int) (*model.Room, error) {
	return gr.repo.FindRoom(ctx, map[string]interface{}{"id": id})
}