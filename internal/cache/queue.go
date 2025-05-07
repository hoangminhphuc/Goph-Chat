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

	pipe := mq.rdb.Pipeline() 
	pipe.LPush(ctx, recentKey, payload) // push on head of the list
	pipe.LTrim(ctx, recentKey, 0, recentCacheSize-1) // cap at recentCacheSize
	pipe.Expire(ctx, recentKey, recentTTL) // expire after TTL
	pipe.RPush(ctx, queueKey, payload)

	_, err := pipe.Exec(ctx)
	return err
}
