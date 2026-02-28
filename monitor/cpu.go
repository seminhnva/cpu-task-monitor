package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUMonitor struct {
}

func (m *CPUMonitor) Name() string {
	return "CPU"
}
func (m *CPUMonitor) GetUsage(ctx context.Context) string {
	cpuStat, err := cpu.PercentWithContext(ctx, 1*time.Second, false)
	if err != nil && len(cpuStat) == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%.2f%%", cpuStat[0])
}
