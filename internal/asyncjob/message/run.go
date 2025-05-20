package message

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/internal/asyncjob"
	"github.com/hoangminhphuc/goph-chat/module/message/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RunBackgroundWorkers(serviceCtx serviceHub.ServiceHub) {
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
	job1 := ProcessMessagesToDB(rdb, msgSaver)
	job2 := InitMessageIdCounter(rdb, db)

	messageJobGroups := asyncjob.NewJobGroup(job1, job2)
	go func() {
		if err := messageJobGroups.Run(ctx); err != nil && err != context.Canceled {
			log.Printf("MessageJobGroups terminated with error: %v", err)
		}
	}()
}