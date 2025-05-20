package message

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hoangminhphuc/goph-chat/internal/asyncjob"
	"github.com/hoangminhphuc/goph-chat/internal/cache"
	"github.com/hoangminhphuc/goph-chat/module/message/model"
	"github.com/redis/go-redis/v9"
)

// ! WRITE-BEHIND PATTERN

// ! RunWorker persists messages from Redis queue into DB
func ProcessMessagesToDB(rdb *redis.Client, saver cache.MessageSaver) asyncjob.Job {
  return asyncjob.NewJob("run-worker", func(ctx context.Context) error {
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
// func RunMessageWorker(serviceCtx serviceHub.ServiceHub) {
//  ctx, cancel := context.WithCancel(context.Background())
//  sigs := make(chan os.Signal, 1)
//  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//  go func() {
//      <-sigs
//      cancel()
//  }()


//  rdb := serviceCtx.MustGetService(common.PluginRedisMain).(*redis.Client)
//  db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)
//  msgSaver := repository.NewSQLRepo(db)


//  // build and start the worker job
//  job := ProcessMessagesToDB(rdb, msgSaver)
//  go func() {
//      if err := job.Execute(ctx); err != nil && err != context.Canceled {
//          log.Printf("RunWorkerJob terminated with error: %v", err)
//      }
//  }()
// }


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

  msgKey := fmt.Sprintf("message:%s", res[1])
  lockKey := fmt.Sprintf("lock:%s", msgKey)
  locked, _ := rdb.Exists(ctx, lockKey).Result()

  if locked == 1 {
    // message is being updated, skip or retry later
    log.Printf("message %s is locked, skipping", msgKey)
    return nil
  }

  data, err := rdb.HGetAll(ctx, msgKey).Result()
  if err != nil {
    return err
  }

  var msg model.Message
	if err := json.Unmarshal([]byte(data["payload"]), &msg); err != nil {
		return fmt.Errorf("json unmarshal error: %w", err)
	}

  return saver.SaveMessage(ctx, &msg)
}
