package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"

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

func GetTopProcesses(ctx context.Context) string {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return fmt.Sprintf("[GetTopProcesses] Could not retrieve memory info :%v\n ", err)
	}
	totalMemory := vmStat.Total
	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return fmt.Sprintf("[GetTopProcesses] Could not retrieve process info :%v\n ", err)
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, proc := range processes {
		wg.Add(1)
		go func(proc *process.Process) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				name, err := proc.NameWithContext(ctx)
				if err != nil {
					return
				}
				cpuPercent, err := proc.CPUPercentWithContext(ctx)
				if err != nil {
					return
				}
				//RSS => Resident Set Size, the non-swapped physical memory a process has used.
				//VMS => Virtual Memory Size, the total amount of virtual memory used by the process.
				memInfo, err := proc.MemoryInfoWithContext(ctx)
				if err != nil {
					return
				}
				ramPercent := float64(memInfo.RSS) / float64(totalMemory) * 100

				createTime, ererr := proc.CreateTimeWithContext(ctx) // miliseconds
				if ererr != nil {
					return
				}
				runningTime := time.Since(time.Unix(createTime/1000, 0))
				if cpuPercent > 5 || ramPercent > 5 {
					mu.Lock()
					procStat := models.ProcStat{
						ID:          proc.Pid,
						Name:        name,
						CPU:         cpuPercent,
						Memory:      memInfo.RSS,
						RamPercent:  ramPercent,
						RunningTime: runningTime,
					}
					fmt.Printf("PID: %d, Name: %s, CPU: %.2f%%, Memory: %d bytes (%.2f%%), Running Time: %s\n", procStat.ID, procStat.Name, procStat.CPU, procStat.Memory, procStat.RamPercent, procStat.RunningTime)
					fmt.Println("==============")
					mu.Unlock()
				}
			}

		}(proc)
	}
	wg.Wait()
	return "processes"
}
