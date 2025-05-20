package business

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/internal/server/websocket"
	"github.com/hoangminhphuc/goph-chat/module/message/dto"
)

type UpdateMessageSQLRepo interface {
	GetMessageByID(ctx context.Context, msgID int) (*websocket.Message, error)
	UpdateMessageByID(ctx context.Context, msgID int, data *dto.MessageUpdate) error
}

type UpdateMessageRedisRepo interface {
	GetMessageInRecent(ctx context.Context, roomID, userID, msgID int) (bool, error)
	UpdateHashMessage(ctx context.Context, msgID int, msgData *dto.MessageUpdate) error
}

type UpdateMessageBusiness struct {
	sqlRepo UpdateMessageSQLRepo
	redisRepo UpdateMessageRedisRepo
}

func NewUpdateMessageBusiness(sqlRepo UpdateMessageSQLRepo, 
	redisRepo UpdateMessageRedisRepo) *UpdateMessageBusiness {
	return &UpdateMessageBusiness{
		sqlRepo: sqlRepo,
		redisRepo: redisRepo,
	}
}

func (b *UpdateMessageBusiness) UpdateMessageByID(
    ctx context.Context, msgID int, data *dto.MessageUpdate) error {
			
    // Check if message is in recent list
    recent, err := b.redisRepo.GetMessageInRecent(ctx, data.RoomID, data.UserID, msgID)
    if err != nil {
        return err
    }

    // Check if message is in recent list
    record, err := b.sqlRepo.GetMessageByID(ctx, msgID)
    if err != nil {
        return err
    }

		// If record not found, update only in Redis
    if record == nil {
        if recent {
            return b.redisRepo.UpdateHashMessage(ctx, msgID, data)
        }
        return nil
    }

    // If record found, update in DB
    if err := b.sqlRepo.UpdateMessageByID(ctx, msgID, data); err != nil {
        return err
    }

    // Then update in Redis
    if recent {
        if err := b.redisRepo.UpdateHashMessage(ctx, msgID, data); err != nil {
            return err
        }
    }

    return nil
}



