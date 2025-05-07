package asyncjob

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/internal/cache"
	"github.com/hoangminhphuc/goph-chat/module/message/model"
	"github.com/hoangminhphuc/goph-chat/module/message/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ! WRITE-BEHIND PATTERN

// ! RunWorker persists messages from Redis queue into DB
func NewRunWorkerJob(rdb *redis.Client, saver cache.MessageSaver) Job {
  return NewJob("run-worker", func(ctx context.Context) error {
    for {
      select {
      case <-ctx.Done():
        log.Println("RunWorkerJob context canceled, exiting worker loop")
        return ctx.Err()
      default:
        err := processQueues(ctx, rdb, saver)
        if err != nil && ctx.Err() == nil {
          log.Printf("worker process error: %v", err)
        }
      }
    }
  })
}


func RunMessageWorker(serviceCtx serviceHub.ServiceHub) {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
			<-sigs
			cancel()
	}()

	rdb := serviceCtx.MustGetService(common.PluginRedisMain).(*redis.Client)
	db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)
	msgSaver := repository.NewSQLRepo(db)

	// build and start the worker job
	job := NewRunWorkerJob(rdb, msgSaver)
	go func() {
			if err := job.Execute(ctx); err != nil && err != context.Canceled {
					log.Printf("RunWorkerJob terminated with error: %v", err)
			}
	}()
}

func processQueues(ctx context.Context, rdb *redis.Client, saver cache.MessageSaver) error {
	var cursor uint64
	pattern := "queues:room:*:*"

	for {
		keys, next, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		cursor = next

		for _, key := range keys {
			if err := popAndSave(ctx, rdb, key, saver); err != nil {
				log.Printf("error processing %s: %v", key, err)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func popAndSave(ctx context.Context, rdb *redis.Client, key string, saver cache.MessageSaver) error {
	res, err := rdb.BLPop(ctx, 0, key).Result()
	if err != nil || len(res) < 2 {
		return err
	}

	var msg model.Message
	if err := json.Unmarshal([]byte(res[1]), &msg); err != nil {
		return fmt.Errorf("json unmarshal error: %w", err)
	}

	return saver.SaveMessage(ctx, &msg)
}

// func parseKey(key string) (room, user string, ok bool) {
// 	parts := strings.Split(key, ":")
// 	if len(parts) != 4 {
// 		return "", "", false
// 	}
// 	return parts[2], parts[3], true
// }
