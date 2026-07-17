package loaders

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type CPULoadStatus struct {
	Active bool
	Cores  int
	EndsAt time.Time
}

type MemoryLoadStatus struct {
	Active   bool
	AmountMB int
}

var mu sync.Mutex
var MemoryHolder []byte
var cpuLoad CPULoadStatus
var memoryLoad MemoryLoadStatus
var cpuLoadGeneration uint64

func StartCPULoad(cores int, timeout int) {
	mu.Lock()
	cpuLoadGeneration++
	generation := cpuLoadGeneration
	cpuLoad = CPULoadStatus{
		Active: true,
		Cores:  cores,
		EndsAt: time.Now().Add(time.Second * time.Duration(timeout)),
	}
	mu.Unlock()

	for i := 0; i < cores; i++ {
		logrus.Info(fmt.Sprintf("Starting CPU load on %d cores for %d seconds", cores, timeout))
		go func() {
			timer := time.NewTimer(time.Second * time.Duration(timeout))
			for {
				//nolint:staticcheck
				select {
				case <-timer.C:
					logrus.Info("CPU load timeout - back to normal...")
					return
				default:
				}
			}
		}()
	}

	time.AfterFunc(time.Second*time.Duration(timeout), func() {
		mu.Lock()
		defer mu.Unlock()
		if cpuLoadGeneration == generation {
			cpuLoad.Active = false
		}
	})
}

// CPUStatus reports the most recently requested synthetic CPU load.
func CPUStatus() CPULoadStatus {
	mu.Lock()
	defer mu.Unlock()

	return cpuLoad
}

func StartMemLeak(amountMb int) {
	mu.Lock()
	memoryLoad = MemoryLoadStatus{Active: amountMb > 0, AmountMB: amountMb}
	mu.Unlock()

	go func() {
		holder := make([]byte, amountMb*1024*1024)

		mu.Lock()
		MemoryHolder = holder
		mu.Unlock()
	}()
}

// MemoryStatus reports the most recently requested synthetic memory footprint.
func MemoryStatus() MemoryLoadStatus {
	mu.Lock()
	defer mu.Unlock()

	return memoryLoad
}
