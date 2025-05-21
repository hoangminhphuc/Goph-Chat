package business

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/module/message/model"
)

type GetRecentMessagesRepo interface {
	GetRecentMessages(ctx context.Context, roomID, userID int) ([]model.Message, error)
}

type GetRecentMessagesBusiness struct {
	repo GetRecentMessagesRepo
}

func NewGetRecentMessagesBusiness(repo GetRecentMessagesRepo) *GetRecentMessagesBusiness {
	return &GetRecentMessagesBusiness{repo: repo}
}

func (b *GetRecentMessagesBusiness) GetRecentMessages(ctx context.Context, 
	roomID, userID int) ([]model.Message, error) {
		return b.repo.GetRecentMessages(ctx, roomID, userID)
}