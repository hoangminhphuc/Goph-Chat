package cache


import (
  "context"
  "fmt"
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
func (mq *MessageQueue) CacheAndQueue(ctx context.Context, room string, userId string, payload []byte) error {
  recentKey := fmt.Sprintf("messages:room:%s:user:%s:recent", room, userId)
  queueKey  := fmt.Sprintf("queues:room:%s:user:%s", room, userId)

  msgID, err := mq.rdb.Get(ctx, "next_message_id").Result()
  if err != nil {
      return err
  }
  msgKey := fmt.Sprintf("message:%s", msgID)

  pipe := mq.rdb.Pipeline()
  pipe.Incr(ctx, "next_message_id") // increment message ID
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

  return err
}
