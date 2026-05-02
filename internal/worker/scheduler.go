package worker

import (
	"context"
	"log"
	"time"

	"blog-api/internal/service"
)

type Scheduler struct {
	postService *service.PostService
	interval    time.Duration
}

func NewScheduler(postService *service.PostService, interval time.Duration) *Scheduler {
	return &Scheduler{
		postService: postService,
		interval:    interval,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("scheduler stopped")
				return
			case <-ticker.C:
				if err := s.postService.PublishScheduled(context.Background(), time.Now()); err != nil {
					log.Printf("scheduler publish error: %v", err)
				}
			}
		}
	}()
}
