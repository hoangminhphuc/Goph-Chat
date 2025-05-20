package message

import (
	"context"
	"log"

	// serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/internal/asyncjob"
	"github.com/hoangminhphuc/goph-chat/module/message/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)


func InitMessageIdCounter(rdb *redis.Client, db *gorm.DB) asyncjob.Job {
  return asyncjob.NewJob("run-worker", func(ctx context.Context) error {
    err := changeList(ctx, rdb, db)
    if err != nil && ctx.Err() == nil {
      log.Printf("worker process error: %v", err)
    }
    return nil
  })
}


func changeList(ctx context.Context, rdb *redis.Client, db *gorm.DB) error {
  repo := repository.NewSQLRepo(db)

  id, err := repo.GetLatestID(ctx)
  if err != nil {
    return err
  }

  _, err = rdb.Set(ctx, "next_message_id", id, 0).Result()
  if err != nil {
      return err
  }
	
  return nil
}
