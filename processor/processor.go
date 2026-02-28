package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/seminhnva/cpu-task-monitor/models"
)

func RunMonitor(ctx context.Context, wg *sync.WaitGroup, statCh chan models.SystemStat, monitor models.Monitor) {
	defer wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Monitor stopped: %s\n", monitor.Name())
			return
		case <-ticker.C:
			statCh <- models.SystemStat{
				Name:  monitor.Name(),
				Value: monitor.GetUsage(ctx),
			}
		}
	}
}
