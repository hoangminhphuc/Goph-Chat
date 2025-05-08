package business

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/internal/server/websocket"
)

type GetRecentMessagesRepo interface {
	GetRecentMessages(ctx context.Context, roomID, userID int) ([]websocket.Message, error)
}

type GetRecentMessagesBusiness struct {
	repo GetRecentMessagesRepo
}

func NewGetRecentMessagesBusiness(repo GetRecentMessagesRepo) *GetRecentMessagesBusiness {
	return &GetRecentMessagesBusiness{repo: repo}
}

func (b *GetRecentMessagesBusiness) GetRecentMessages(ctx context.Context, 
	roomID, userID int) ([]websocket.Message, error) {
		return b.repo.GetRecentMessages(ctx, roomID, userID)
}