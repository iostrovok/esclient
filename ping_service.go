package esclient

import (
	"context"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
)

type pingService struct {
	mc sync.RWMutex

	sniffDuration time.Duration
	reConnect     func() error
	printf        func(format string, v ...interface{})

	pingService []*elastic.PingService
}

func newPingService(sniffDuration time.Duration, reConnect func() error, printf func(format string, v ...interface{})) *pingService {
	return &pingService{
		pingService:   []*elastic.PingService{},
		sniffDuration: sniffDuration,
		reConnect:     reConnect,
		printf:        printf,
	}
}

func (ps *pingService) Add(e *elastic.PingService) {
	ps.pingService = append(ps.pingService, e)
}

func (ps *pingService) runSniff(ctx context.Context) {
	// sleep before first attempt
	time.Sleep(ps.sniffDuration)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(ps.sniffDuration):
			if !ps.checkConnections(ctx) {
				ps.reConnect()
			}
		}
	}
}

func (ps *pingService) checkConnections(ctx context.Context) bool {
	ps.mc.RLock()
	defer ps.mc.RUnlock()

	for i := range ps.pingService {
		_, statusCode, err := ps.pingService[i].Do(ctx)
		if err != nil || statusCode < 200 || statusCode >= 300 {
			ps.printf("Reconnect statusCode: %d. err: %v. res: %v", statusCode, err)
			return false
		}
	}

	return true
}
