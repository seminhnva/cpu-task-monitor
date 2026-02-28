package models

import (
	"context"
	"sync"
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
