package loaders

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_StartCPULoadShouldReportActiveStatusImmediately(t *testing.T) {
	StartCPULoad(2, 5)

	status := CPUStatus()

	assert.True(t, status.Active)
	assert.Equal(t, 2, status.Cores)
	assert.True(t, status.EndsAt.After(time.Now()))
}

func Test_CPUStatusShouldBecomeInactiveAfterTimeoutElapses(t *testing.T) {
	StartCPULoad(1, 1)

	assert.True(t, CPUStatus().Active)

	time.Sleep(time.Millisecond * 1200)

	assert.False(t, CPUStatus().Active)
}

func Test_StartMemLeakShouldReportRequestedFootprintImmediately(t *testing.T) {
	StartMemLeak(5)

	status := MemoryStatus()

	assert.True(t, status.Active)
	assert.Equal(t, 5, status.AmountMB)
}

func Test_StartMemLeakWithZeroShouldReportInactiveStatus(t *testing.T) {
	StartMemLeak(7)
	assert.True(t, MemoryStatus().Active)

	StartMemLeak(0)

	status := MemoryStatus()
	assert.False(t, status.Active)
	assert.Equal(t, 0, status.AmountMB)
}
