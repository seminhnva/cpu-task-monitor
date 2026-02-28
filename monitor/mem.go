package monitor

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryMonitor struct {
}

func (m *MemoryMonitor) Name() string {
	return "Memory"
}

func (m *MemoryMonitor) GetUsage(ctx context.Context) string {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return "N/A"
	}
	return fmt.Sprintf("%.2f%%", vmStat.UsedPercent)
}
