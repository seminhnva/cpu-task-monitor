package monitor

import (
	"context"
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskMonitor struct {
}

func (m *DiskMonitor) Name() string {
	return "Disk"
}
func (m *DiskMonitor) GetUsage(ctx context.Context) string {
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:\\"
	}
	diskStat, err := disk.UsageWithContext(ctx, path)
	if err != nil {
		return fmt.Sprintf("[Disk Monitor] Could not retrieve disk info :%v\n ", err)
	}
	return fmt.Sprintf("%.2f%%", diskStat.UsedPercent)
}
