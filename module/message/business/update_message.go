package business

import (
	"context"

	"github.com/hoangminhphuc/goph-chat/module/message/dto"
	"github.com/hoangminhphuc/goph-chat/module/message/model"
)

type UpdateMessageSQLRepo interface {
	GetMessageByID(ctx context.Context, msgID int) (*model.Message, error)
	UpdateMessageByID(ctx context.Context, msgID int, data *dto.MessageUpdate) error
}

type UpdateMessageRedisRepo interface {
	GetMessageInRecent(ctx context.Context, roomID, userID, msgID int) (bool, error)
	UpdateHashMessage(ctx context.Context, msgID int, msgData *dto.MessageUpdate) (*model.Message, error)
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
    ctx context.Context, msgID int, data *dto.MessageUpdate) (*model.Message, error) {
			
    // Check if message is in recent list
    recent, err := b.redisRepo.GetMessageInRecent(ctx, data.RoomID, data.UserID, msgID)
    if err != nil {
        return nil, err
    }

    // Check if message is in DB
    record, err := b.sqlRepo.GetMessageByID(ctx, msgID)
    if err != nil {
        return nil, err
    }

	// If record not found, update only in Redis
    if record == nil {
        if recent { // Return message in Redis not in DB
            return b.redisRepo.UpdateHashMessage(ctx, msgID, data)
        }
        return nil, nil
    }

    // If record found, update in DB
    if err := b.sqlRepo.UpdateMessageByID(ctx, msgID, data); err != nil {
        return nil, err
    }

    updated, err := b.sqlRepo.GetMessageByID(ctx, msgID)
    if err != nil {
        return nil, err
    }

    // Then update in Redis
    if recent {
        if _, err := b.redisRepo.UpdateHashMessage(ctx, msgID, data); err != nil {
            return nil, err
        }
    }
    // Return message in DB
    return updated, nil
}



