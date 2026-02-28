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
		return fmt.Sprintf("[Memory Monitor] Could not retrieve memory info :%v\n ", err)
	}
	return fmt.Sprintf("%.2f%%", vmStat.UsedPercent)
}
