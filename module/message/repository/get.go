package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hoangminhphuc/goph-chat/internal/server/websocket"
	"github.com/hoangminhphuc/goph-chat/module/message/model"
	"github.com/redis/go-redis/v9"
)

// ! REDIS REPO

func (r *redisRepo) GetRecentMessages(ctx context.Context, roomID, userID int) ([]websocket.Message, error) {
	key := fmt.Sprintf("messages:room:%d:user:%d:recent", roomID, userID)

	vals, err := r.rdb.LRange(ctx, key, 0, -1).Result() 
	if err != nil {
		return nil, err
	}

	pipe := r.rdb.Pipeline()
	for _, id := range vals {
		pipe.HGetAll(ctx, fmt.Sprintf("message:%s", id))
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	var msgs []websocket.Message
	for i := len(cmds) - 1; i >= 0; i-- { 
		// Each cmd is a *redis.MapStringStringCmd
		data := cmds[i].(*redis.MapStringStringCmd).Val()
		var m websocket.Message

		if err := json.Unmarshal([]byte(data["payload"]), &m); err != nil {
			return nil, err
		}

		m.ChatUser, m.RoomID = 0, 0 // Just showing the message content
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (r *redisRepo) GetMessageInRecent(ctx context.Context, roomID, userID, msgID int) (bool, error) {
	key := fmt.Sprintf("messages:room:%d:user:%d:recent", roomID, userID)
	
	vals, err := r.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return false, err
	}

	target := strconv.Itoa(msgID)
	for _, val := range vals {
		if val == target {
			return true, nil
		}
	}
	return false, nil
}


// ! SQL REPO
func (s *sqlRepo) GetMessageByID(ctx context.Context, msgID int) (*websocket.Message, error) {
	var msg websocket.Message
	db := s.db.Table(model.Message{}.TableName())
	err := db.Where("id = ?", msgID).First(&msg).Error
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (s *sqlRepo) GetLatestID(ctx context.Context) (int, error) {
	var lastID int
	db := s.db.Table(model.Message{}.TableName())
	err := db.Select("COALESCE(MAX(id), 0)").Scan(&lastID).Error
	if err != nil {
		return 0, err
	}

	return lastID, nil
}