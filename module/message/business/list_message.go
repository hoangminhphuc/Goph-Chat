package business

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/message/model"
)

type ListMessageRepo interface {
	ListMessage(ctx context.Context, roomID int, paging *utils.Paging) ([]model.Message, error)
}

type ListMessageBusiness struct {
	repo ListMessageRepo
}

func NewListMessageBusiness(repo ListMessageRepo) *ListMessageBusiness {
	return &ListMessageBusiness{repo: repo}
}

func (b *ListMessageBusiness) ListMessage(ctx context.Context, roomID int, paging *utils.Paging) ([]model.Message, error) {
	return b.repo.ListMessage(ctx, roomID, paging)
}
