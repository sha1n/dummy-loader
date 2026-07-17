package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sha1n/dummy-loader/server/loaders"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_HandleStatusRequestShouldReportCurrentLoad(t *testing.T) {
	loaders.StartCPULoad(2, 5)
	loaders.StartMemLeak(10)

	c, w := testContext()
	HandleStatusRequest(c)

	assert.Equal(t, 200, w.Code)

	var body statusResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.True(t, body.CPU.Active)
	assert.Equal(t, 2, body.CPU.Cores)
	assert.True(t, body.CPU.RemainingSec > 0)
	assert.True(t, body.Memory.Active)
	assert.Equal(t, 10, body.Memory.AmountMB)
}

func Test_HandleMetricsRequestShouldSetPrometheusContentType(t *testing.T) {
	c, w := testContext()
	HandleMetricsRequest(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "text/plain; version=0.0.4; charset=utf-8", w.Header().Get("Content-Type"))
}

func Test_RenderMetricsShouldFormatActiveGauges(t *testing.T) {
	cpu := loaders.CPULoadStatus{Active: true, Cores: 4, EndsAt: time.Now().Add(10 * time.Second)}
	mem := loaders.MemoryLoadStatus{Active: true, AmountMB: 256}

	body := renderMetrics(cpu, mem)

	assert.Contains(t, body, "# TYPE dummy_loader_cpu_load_active gauge")
	assert.Contains(t, body, "dummy_loader_cpu_load_active 1")
	assert.Contains(t, body, "dummy_loader_cpu_load_cores 4")
	assert.Contains(t, body, "dummy_loader_memory_active 1")
	assert.Contains(t, body, "dummy_loader_memory_footprint_mb 256")
}

func Test_RenderMetricsShouldReportZeroGaugesWhenInactive(t *testing.T) {
	cpu := loaders.CPULoadStatus{Active: false}
	mem := loaders.MemoryLoadStatus{Active: false}

	body := renderMetrics(cpu, mem)

	assert.Contains(t, body, "dummy_loader_cpu_load_active 0")
	assert.Contains(t, body, "dummy_loader_cpu_load_remaining_seconds 0")
	assert.Contains(t, body, "dummy_loader_memory_active 0")
	assert.Contains(t, body, "dummy_loader_memory_footprint_mb 0")
}

func testContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}
