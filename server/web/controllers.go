package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sha1n/dummy-loader/server/loaders"
	"github.com/sirupsen/logrus"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var startupTime = time.Now()

func HandleHealthCheck(c *gin.Context) {
	c.Status(200)
}

func HandleReadinessCheck(c *gin.Context) {
	// getting ready only after 20 seconds
	if time.Now().Unix() >= (startupTime.Add(time.Second * 20).Unix()) {
		c.Status(200)
	} else {
		logrus.Warn("Server not ready yet...")
		c.Status(503)
	}
}

func HandleCPULoadRequest(c *gin.Context) {
	timeSec := c.Request.URL.Query().Get("time-sec")
	coresStr := c.Request.URL.Query().Get("cores")
	var cores = runtime.NumCPU()
	if coresStr != "" {
		requestedCores, e := strconv.Atoi(coresStr)
		if e != nil {
			c.Status(400)
		}
		cores = int(math.Min(float64(requestedCores), float64(runtime.NumCPU())))
	}
	if timeSec != "" {
		i, e := strconv.Atoi(timeSec)
		if e != nil {
			c.Status(400)
		} else {
			c.Status(202)
			loaders.StartCPULoad(cores, i)
			_, _ = c.Writer.WriteString(fmt.Sprintf("CPU load started for %s seconds on %d cores", timeSec, cores))
		}

	} else {
		c.Redirect(307, "/docs")
	}
}

func HandleMemLeakRequest(c *gin.Context) {
	amountMb := c.Request.URL.Query().Get("amount-mb")
	if amountMb != "" {
		c.Status(202)
		i, e := strconv.Atoi(amountMb)
		if e != nil {
			c.Status(400)
		} else {
			loaders.StartMemLeak(i)
			c.Status(202)
			_, _ = c.Writer.WriteString(fmt.Sprintf("Memory footprint set to %sMB", amountMb))
		}
	} else {
		c.Redirect(307, "/docs")
	}
}

type cpuStatusResponse struct {
	Active       bool `json:"active"`
	Cores        int  `json:"cores"`
	RemainingSec int  `json:"remainingSec"`
}

type memoryStatusResponse struct {
	Active   bool `json:"active"`
	AmountMB int  `json:"amountMb"`
}

type statusResponse struct {
	CPU    cpuStatusResponse    `json:"cpu"`
	Memory memoryStatusResponse `json:"memory"`
}

func HandleStatusRequest(c *gin.Context) {
	cpu := loaders.CPUStatus()
	mem := loaders.MemoryStatus()

	c.JSON(200, statusResponse{
		CPU: cpuStatusResponse{
			Active:       cpu.Active,
			Cores:        cpu.Cores,
			RemainingSec: remainingSeconds(cpu),
		},
		Memory: memoryStatusResponse{
			Active:   mem.Active,
			AmountMB: mem.AmountMB,
		},
	})
}

func HandleMetricsRequest(c *gin.Context) {
	body := renderMetrics(loaders.CPUStatus(), loaders.MemoryStatus())
	c.Data(200, "text/plain; version=0.0.4; charset=utf-8", []byte(body))
}

func remainingSeconds(cpu loaders.CPULoadStatus) int {
	if !cpu.Active {
		return 0
	}

	remaining := int(math.Ceil(time.Until(cpu.EndsAt).Seconds()))
	if remaining < 0 {
		return 0
	}

	return remaining
}

func renderMetrics(cpu loaders.CPULoadStatus, mem loaders.MemoryLoadStatus) string {
	var b strings.Builder

	writeGauge(&b, "dummy_loader_cpu_load_active", "Whether synthetic CPU load is currently active (1) or not (0).", boolToInt(cpu.Active))
	writeGauge(&b, "dummy_loader_cpu_load_cores", "Number of CPU cores currently under synthetic load.", cpu.Cores)
	writeGauge(&b, "dummy_loader_cpu_load_remaining_seconds", "Seconds remaining until synthetic CPU load stops.", remainingSeconds(cpu))
	writeGauge(&b, "dummy_loader_memory_active", "Whether a synthetic memory footprint is currently active (1) or not (0).", boolToInt(mem.Active))
	writeGauge(&b, "dummy_loader_memory_footprint_mb", "Size in megabytes of the synthetic memory footprint currently held.", mem.AmountMB)

	return b.String()
}

func writeGauge(b *strings.Builder, name string, help string, value int) {
	fmt.Fprintf(b, "# HELP %s %s\n# TYPE %s gauge\n%s %d\n", name, help, name, name, value)
}

func boolToInt(v bool) int {
	if v {
		return 1
	}

	return 0
}
