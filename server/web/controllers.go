package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sha1n/dummy-loader/server/loaders"
	"github.com/sirupsen/logrus"
	"math"
	"runtime"
	"strconv"
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

func HandleCpuLoadRequest(c *gin.Context) {
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
			loaders.StartCpuLoad(cores, i)
			c.Writer.WriteString(fmt.Sprintf("CPU load started for %s seconds on %d cores", timeSec, cores))
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
			c.Writer.WriteString(fmt.Sprintf("Memory footprint set to %sMB", amountMb))
		}
	} else {
		c.Redirect(307, "/docs")
	}
}
