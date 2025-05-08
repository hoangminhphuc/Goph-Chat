package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hoangminhphuc/goph-chat/internal/server/websocket"
)

func (r *redisRepo) GetRecentMessages(ctx context.Context, roomID, userID int) ([]websocket.Message, error) {
	key := fmt.Sprintf("messages:room:%d:user:%d:recent", roomID, userID)

	vals, err := r.rdb.LRange(ctx, key, 0, -1).Result() 
	if err != nil {
		return nil, err
	}

	var msgs []websocket.Message
	for i := len(vals) - 1; i >= 0; i-- { 
		var m websocket.Message
		if err := json.Unmarshal([]byte(vals[i]), &m); err != nil {
			// skip or handle bad JSON
			continue
		}
		m.ChatUser, m.RoomID = 0, 0 // Just showing the message content
		msgs = append(msgs, m)
	}
	return msgs, nil
}