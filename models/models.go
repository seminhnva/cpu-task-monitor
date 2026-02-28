package models

import (
	"context"
	"sync"
	"time"
)

type SystemStat struct {
	Name  string
	Value string
}

type Monitor interface {
	Name() string
	GetUsage(ctx context.Context) string
}

var (
	Stats     = map[string]SystemStat{}
	StatMutex sync.Mutex
)

type ProcStat struct {
	ID          int32
	Name        string
	CPU         float64
	Memory      uint64
	RamPercent  float64
	RunningTime time.Duration
}
