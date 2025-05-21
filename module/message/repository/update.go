package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hoangminhphuc/goph-chat/module/message/dto"
	"github.com/hoangminhphuc/goph-chat/module/message/model"
)

const lockTTL = 10 * time.Second

// ! SQL REPO
func (s *sqlRepo) UpdateMessageByID(ctx context.Context, msgID int, data *dto.MessageUpdate) error {
	if err :=s.db.Table(data.TableName()).Where("id = ?", msgID).
	Updates(data).Error; err != nil {
		return err
	}
	return nil
}



// ! REDIS REPO
// Update in hash redis
func (r *redisRepo) UpdateHashMessage(ctx context.Context, 
		msgID int, msgData *dto.MessageUpdate) (*model.Message, error) {

	msgKey := fmt.Sprintf("message:%d", msgID)
	lockKey := "lock:" + msgKey
	acquired, err := r.rdb.SetNX(ctx, lockKey, 1, lockTTL).Result()

	if err != nil {
		return nil, err
	}

	if !acquired { // If this gets called by another process, we should wait
		return nil, fmt.Errorf("resource is locked, try again later")
	}
	defer r.rdb.Del(ctx, lockKey) // only safe because we know no one else could've set it



	data, err := r.rdb.HGetAll(ctx, msgKey).Result()
  if err != nil {
    return nil, err
  }

  var msg model.Message
	if err := json.Unmarshal([]byte(data["payload"]), &msg); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	msg.Content = msgData.Content
	updatedPayload, _ := json.Marshal(msg)

	if err := r.rdb.HSet(ctx, msgKey, map[string]interface{}{
		"payload": []byte(updatedPayload),
	}).Err();err != nil {
		return nil, err
	}

	return &msg, nil
}