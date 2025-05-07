package business

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/room/model"
)

type ListRoomRepo interface {
	ListRoom(ctx context.Context, paging *utils.Paging) ([]model.Room, error)
}

type ListRoomBusiness struct {
	repo ListRoomRepo
}

func NewListRoomBusiness(repo ListRoomRepo) *ListRoomBusiness {
	return &ListRoomBusiness{repo: repo}
}

func (lr *ListRoomBusiness) ListRoom(ctx context.Context, paging *utils.Paging) ([]model.Room, error) {
	return lr.repo.ListRoom(ctx, paging)
}