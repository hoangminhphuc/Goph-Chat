package business

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/module/room/dto"
)

type CreateRoomRepo interface {
	CreateRoom(ctx context.Context, data *dto.RoomCreation) error
	
}

type CreateRoomBusiness struct {
	repo CreateRoomRepo
}

func NewCreateRoomBusiness(repo CreateRoomRepo) *CreateRoomBusiness {
	return &CreateRoomBusiness{repo: repo}
}

func (cr *CreateRoomBusiness) CreateRoom(ctx context.Context, data *dto.RoomCreation) error {
	if err := cr.repo.CreateRoom(ctx, data); err != nil {
		common.ErrCannotCreate("room", err)
	}

	return nil
}
