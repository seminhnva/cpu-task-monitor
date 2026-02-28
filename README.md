# cpu-task-monitor

A real-time system monitor written in Go. Tracks CPU, memory, network, and disk — and reports the top resource-hungry processes every few seconds.

## Demo

```
$ go run main.go

=== System status ===
[CPU]: 34.20%
[Memory]: 61.05%
[Network]: Send: 1024 KB, Recv: 4096 KB
[Disk]: 55.30%

====== Top CPU consuming processes =====
1.[1234] chrome    CPU: 18.50%, RAM: 512.30 MB( 6.25%), Running Time: 2h15m0s
2.[5678] gopls     CPU:  5.20%, RAM: 210.10 MB( 2.56%), Running Time: 45m10s
3.[9012] node      CPU:  3.10%, RAM: 180.40 MB( 2.20%), Running Time: 1h5m0s

====Top RAM consuming processes ===
1.[1234] chrome    CPU: 18.50%, RAM: 512.30 MB( 6.25%), Running Time: 2h15m0s
2.[3456] slack     CPU:  1.20%, RAM: 430.00 MB( 5.25%), Running Time: 3h0m0s
3.[5678] gopls     CPU:  5.20%, RAM: 210.10 MB( 2.56%), Running Time: 45m10s

[2025-01-15T10:32:01+07:00] ALERT :Memory = 61.05%
```

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                       main.go                        │
│                                                      │
│  ┌────────────┐  ┌──────────┐  ┌─────┐  ┌──────┐   │
│  │ CPUMonitor │  │MemMonitor│  │ Net │  │ Disk │   │  goroutines
│  └─────┬──────┘  └────┬─────┘  └──┬──┘  └──┬───┘   │
│        └──────────────┴───────────┴─────────┘        │
│                        │ statCh                       │
│                   ┌────▼────┐                         │
│                   │ Stats   │  (map + mutex)           │
│                   └────┬────┘                         │
│              printTicker (3s)                         │
│                   ┌────▼──────────────────────┐       │
│                   │     GetTopProcesses()      │       │
│                   │  worker pool (20 goroutines)│      │
│                   └────┬──────────────────┬────┘      │
│                        │                  │            │
│               ┌────────▼───┐    ┌─────────▼──────┐   │
│               │ alert.log  │    │process_stat.csv │   │
│               └────────────┘    └────────────────-┘   │
└─────────────────────────────────────────────────────┘
```

## Features

- Monitors CPU, Memory, Network, Disk concurrently (sampled every second)
- Prints system stats + top 5 processes by CPU and RAM every 3 seconds
- Logs alerts to `alert.log` when CPU, Memory, or Disk exceeds **60%**
- Exports process snapshots to `process_stat.csv`

## Quick Start

```bash
git clone https://github.com/seminhnva/cpu-task-monitor.git
cd cpu-task-monitor
go mod tidy
go run main.go
```

Runs for 60 seconds then exits cleanly. To change the duration, edit this line in `main.go`:

```go
time.Sleep(60 * time.Second) // change 60 to any number of seconds
```

## Output Files

| File | Description |
|------|-------------|
| `alert.log` | Alert when CPU / Memory / Disk > 60% |
| `process_stat.csv` | Top process snapshot on each tick |

## Project Structure

```
cpu-task-monitor/
├── main.go          # Entry point
├── models/          # Shared types and Monitor interface
├── monitor/         # CPU, Memory, Network, Disk implementations
└── processor/       # RunMonitor loop, process collector, CSV export, alert logger
```

## Adding a New Monitor

Implement the `Monitor` interface:

```go
type Monitor interface {
    Name() string
    GetUsage(ctx context.Context) (value string, alert bool)
}
```

Then register it in `main.go`:

```go
monitorLst := []models.Monitor{
    &monitor.CPUMonitor{},
    &monitor.YourNewMonitor{}, // 👈
}
```

## Dependencies

- [`gopsutil/v4`](https://github.com/shirou/gopsutil) — cross-platform system stats

## License

[MIT](./LICENSE)