package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/seminhnva/cpu-task-monitor/models"
	"github.com/seminhnva/cpu-task-monitor/monitor"
	"github.com/seminhnva/cpu-task-monitor/processor"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitorLst := []models.Monitor{
		&monitor.CPUMonitor{},
		&monitor.MemoryMonitor{},
		&monitor.NetMonitor{},
		&monitor.DiskMonitor{},
	}
	var wg sync.WaitGroup
	statCh := make(chan models.SystemStat)

	for _, monitor := range monitorLst {
		wg.Add(1)
		go processor.RunMonitor(ctx, &wg, statCh, monitor)
	}
	go func() {
		for stat := range statCh {
			models.StatMutex.Lock()
			models.Stats[stat.Name] = stat
			models.StatMutex.Unlock()
		}
	}()
	printTicker := time.NewTicker(3 * time.Second)
	go func() {
		defer printTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-printTicker.C:
				fmt.Println("=== System status ===")
				models.StatMutex.Lock()
				for _, valStat := range models.Stats {
					fmt.Printf("[%s]: %s\n", valStat.Name, valStat.Value)
				}
				models.StatMutex.Unlock()
				fmt.Println(processor.GetTopProcesses(ctx))
			}
		}

	}()
	time.Sleep(60 * time.Second)
	cancel()
	wg.Wait()
	close(statCh)
}
