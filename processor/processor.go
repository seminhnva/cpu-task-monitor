package processor

import (
	"context"
	"fmt"
	"os"
	"sort"
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
	var output string
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return fmt.Sprintf("[GetTopProcesses] Could not retrieve memory info :%v\n ", err)
	}
	totalMemory := vmStat.Total
	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return fmt.Sprintf("[GetTopProcesses] Could not retrieve process info :%v\n ", err)
	}
	var wg sync.WaitGroup
	var cpuList, memList []models.ProcStat
	procChan := make(chan models.ProcStat, len(processes))

	for _, proc := range processes {
		wg.Add(1)
		go func(proc *process.Process) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			name, err := proc.NameWithContext(ctx)
			if err != nil {
				return
			}
			cpuPercent, err := proc.CPUPercentWithContext(ctx)
			if err != nil {
				return
			}
			memInfo, err := proc.MemoryInfoWithContext(ctx)
			if err != nil {
				return
			}
			ramPercent := float64(memInfo.RSS) / float64(totalMemory) * 100

			if cpuPercent < 1 || ramPercent < 1 {
				return
			}

			createTime, ererr := proc.CreateTimeWithContext(ctx) // miliseconds
			if ererr != nil {
				return
			}

			runningTime := time.Since(time.Unix(createTime/1000, 0))
			procStat := models.ProcStat{
				PID:         proc.Pid,
				Name:        name,
				CPU:         cpuPercent,
				Memory:      memInfo.RSS,
				RamPercent:  ramPercent,
				RunningTime: runningTime,
			}
			procChan <- procStat
		}(proc)
	}

	go func() {
		wg.Wait()
		close(procChan)
	}()

	for stat := range procChan {
		if stat.CPU > 1 {
			cpuList = append(cpuList, stat)
		}
		if stat.RamPercent > 1 {
			memList = append(memList, stat)
		}
	}
	sort.Slice(cpuList, func(i, j int) bool {
		return cpuList[i].CPU > cpuList[j].CPU
	})
	sort.Slice(memList, func(i, j int) bool {
		return memList[i].RamPercent > memList[j].RamPercent
	})
	output += "====== Top CPU consuming processes =====\n"
	for i, cpu := range cpuList[:5] {
		output += fmt.Sprintf("%d.[%d] %s CPU: %.2f%%, RAM: %.2f MB( %.2f%%), Running Time: %s\n",
			i+1,
			cpu.PID,
			cpu.Name,
			cpu.CPU,
			float64(cpu.Memory)/(1024*1024),
			cpu.RamPercent,
			cpu.RunningTime)
	}
	output += "====Top RAM consuming processes ===\n"
	for i, mem := range memList[:5] {
		output += fmt.Sprintf("%d.[%d] %s CPU: %.2f%%, RAM: %.2f MB( %.2f%%), Running Time: %s\n",
			i+1,
			mem.PID,
			mem.Name,
			mem.CPU,
			float64(mem.Memory)/(1024*1024),
			mem.RamPercent,
			mem.RunningTime)
	}
	ExportToCsv(cpuList, memList)
	return output
}

func ExportToCsv(cpuList, memList []models.ProcStat) {
	file, err := os.OpenFile("process_stat.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[Export To CSV]sCould not create CSV file: %v\n", err)
		return
	}
	defer file.Close()
	if stat, err := file.Stat(); err == nil && stat.Size() == 0 {
		file.WriteString("Timestamp, PID, Name, CPU%, Memory(MB), RAM%, RunningTime\n")
	}
	timeStamp := time.Now().Format(time.RFC3339)
	for _, cpu := range cpuList[:5] {
		file.WriteString(fmt.Sprintf("%s, %d, %s, %.2f%%, %.2f, %.2f%%, %s\n",
			timeStamp,
			cpu.PID,
			cpu.Name,
			cpu.CPU,
			float64(cpu.Memory)/(1024*1024),
			cpu.RamPercent,
			cpu.RunningTime))
	}
	for _, mem := range memList[:5] {
		file.WriteString(fmt.Sprintf("%s, %d, %s, %.2f%%, %.2f, %.2f%%, %s\n",
			timeStamp,
			mem.PID,
			mem.Name,
			mem.CPU,
			float64(mem.Memory)/(1024*1024),
			mem.RamPercent,
			mem.RunningTime))
	}
}
