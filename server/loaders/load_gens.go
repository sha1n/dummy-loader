package loaders

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

var memoryHolder []byte

func StartCpuLoad(cores int, timeout int) {
	for i := 0; i < cores; i++ {
		logrus.Info(fmt.Sprintf("Starting CPU load on %d cores for %d seconds", cores, timeout))
		go func() {
			timer := time.NewTimer(time.Second * time.Duration(timeout))
			for {
				select {
				case <-timer.C:
					logrus.Info("CPU load timeout - back to normal...")
					return
				default:
				}
			}
		}()
	}
}

func StartMemLeak(amountMb int) {
	go func() {
		memoryHolder = make([]byte, amountMb*1024*1024)
	}()
}
