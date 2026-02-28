package monitor

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/net"
)

type NetMonitor struct {
}

func (m *NetMonitor) Name() string {
	return "Network"
}
func (m *NetMonitor) GetUsage(ctx context.Context) (string, bool) {
	netStat, err := net.IOCountersWithContext(ctx, false)
	if err != nil && len(netStat) == 0 {
		return fmt.Sprintf("[Network Monitor] Could not retrieve network info :%v\n ", err), false
	}
	return fmt.Sprintf("Send: %v KB, Recv: %v KB", netStat[0].BytesSent/1024, netStat[0].BytesRecv/1024), false
}
