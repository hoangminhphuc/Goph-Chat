package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hoangminhphuc/goph-chat/module/message/model"
	"github.com/redis/go-redis/v9"
)


const (
  recentCacheSize = 100
  recentTTL       = 24 * time.Hour
)


type MessageSaver interface {
  SaveMessage(ctx context.Context, message *model.Message) error
}


type MessageQueue struct {
  rdb *redis.Client
}


func NewMessageQueue(rdb *redis.Client) *MessageQueue {
  return &MessageQueue{rdb: rdb}
}


// CacheAndQueue stores message to recent list and async queue
func (mq *MessageQueue) CacheAndQueue(ctx context.Context, room string, userId string, payload []byte) (*model.Message, error) {
  recentKey := fmt.Sprintf("messages:room:%s:user:%s:recent", room, userId)
  queueKey  := fmt.Sprintf("queues:room:%s:user:%s", room, userId)

  mq.rdb.Incr(ctx, "next_message_id") // increment message ID
  msgID, err := mq.rdb.Get(ctx, "next_message_id").Result()
  if err != nil {
      return nil, err
  }

  msgKey := fmt.Sprintf("message:%s", msgID)
  pipe := mq.rdb.Pipeline()
  pipe.LPush(ctx, recentKey, msgID) // push on head of the list
  pipe.LTrim(ctx, recentKey, 0, recentCacheSize-1) // cap at recentCacheSize
  pipe.Expire(ctx, recentKey, recentTTL)

  pipe.HSet(ctx, msgKey, map[string]interface{}{
    "payload": payload, // marshalled message

    // "created": msg.CreatedAt.UnixNano(),
    // "edited":  msg.EditedAt.UnixNano(),
  })
  pipe.Expire(ctx, msgKey, recentTTL)

  pipe.RPush(ctx, queueKey, msgID) // push to tail

  _, err = pipe.Exec(ctx)

  if err != nil {
      return nil, err
  }

  return buildMessage(msgID, payload)
}

func buildMessage(msgID string, payload []byte) (*model.Message, error) {
    var msg model.Message              
    if err := json.Unmarshal(payload, &msg); err != nil {
        return nil, fmt.Errorf("invalid payload JSON: %w", err)
    }

    id, err := strconv.Atoi(msgID)
    if err != nil {
        return nil, fmt.Errorf("invalid message ID: %w", err)
    }

    msg.ID = id                     
    return &msg, nil                   
}