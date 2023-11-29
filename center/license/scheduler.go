package license

import (
	"context"
	"time"
)

type Scheduler struct {
	// key: hash
	frequency string
}

func NewScheduler(frequency string) *Scheduler {
	scheduler := &Scheduler{
		frequency: frequency,
	}

	go scheduler.LoopSyncLicense(context.Background())
	return scheduler
}

func (s *Scheduler) LoopSyncLicense(ctx context.Context) {
	// if s.frequency == "once" {

	// } else if s.frequency == "days" {
	duration := 24 * time.Hour
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			s.syncLicense()
		}
	}
	// }
}

func (s *Scheduler) syncLicense() {
}
