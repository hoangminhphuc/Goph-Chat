package maintenance

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/logger"
	"github.com/hoangminhphuc/goph-chat/internal/asyncjob"
	"github.com/hoangminhphuc/goph-chat/module"
	"gorm.io/gorm"
)

var (
	DefaultCleanupInterval = time.Hour // 1 hour
	DefaultCleanupTTL      = 30 * 24 * time.Hour // 30 days
)

// CleanupService runs periodic cleanup of soft-deleted records via asyncjob
type CleanupService struct {
	Interval  time.Duration // interval between cleanup checks
	TTL       time.Duration // TTL for soft-deleted records
	logger 		logger.ZapLogger
}

func NewCleanupService() *CleanupService {
	return &CleanupService{
		Interval:  DefaultCleanupInterval,
		TTL:       DefaultCleanupTTL,
		logger: logger.NewZapLogger(),
	}
}

// Run starts a ticker and, on each tick, spawns cleanup jobs for all models
func (s *CleanupService) Run(ctx context.Context, serviceCtx serviceHub.ServiceHub) error {
	db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	models := module.GetAllModels()

	for {
		select {
		case <-ctx.Done():
			// service shutdown requested
			return ctx.Err()
		case <-ticker.C:
			// on each tick, build jobs for each model
			var jobs []asyncjob.Job
			expiration := time.Now().Add(-s.TTL)

			for _, m := range models {
				model := m
				table := model.(interface{ TableName() string }).TableName()

				handler := func(ctx context.Context) error {
					// perform deletion for one model
					if err := db.
						Table(table).
						Unscoped().
						Where("deleted_at < ?", expiration).
						Delete(model).
						Error; err != nil {
						return err
					}
					return nil
				}

				jobs = append(jobs, asyncjob.NewJob("cleanup:"+table, handler))
			}

			// run all cleanup jobs in parallel, cancel remaining on first error
			if err := asyncjob.NewJobGroup(jobs...).Run(ctx); err != nil {
				s.logger.Log.Error("Cleanup error:", err)
			}
		}
	}
}

func StartCleanupService(serviceCtx serviceHub.ServiceHub) {
	ctx, cancel := context.WithCancel(context.Background())
	
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
			<-sigs
			cancel()
	}()

	cleanupSvc := NewCleanupService()
	go func() {
			if err := cleanupSvc.Run(ctx, serviceCtx); err != nil && err != context.Canceled {
					cleanupSvc.logger.Log.Error("CleanupService terminated with error:", err)
			}
	}()
}
