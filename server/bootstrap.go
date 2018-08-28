package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sha1n/dummy-loader/server/http"
	"github.com/sha1n/dummy-loader/server/sys"
	"github.com/sha1n/dummy-loader/server/web"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile | log.Ltime | log.LUTC)
}

func main() {
	server := createHttpServer(8080)
	server.StartAsync()

	awaitShutdownSig()
}

func awaitShutdownSig() {
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Waiting for shutdown signal...")

	<-quitChannel
}

func createHttpServer(port int) http.Server {
	server := http.
		NewServer(port).
		WithGetHandler("/health", web.HandleHealthCheck).
		WithGetHandler("/ready", web.HandleReadinessCheck).
		WithGetHandler("/docs", docsHandler).
		WithHtmlTemplates("server/web/help.tpl").
		WithGetHandler("/api/cpu-load", web.HandleCpuLoadRequest).
		WithPostHandler("/api/cpu-load", web.HandleCpuLoadRequest).
		WithGetHandler("/api/mem-footprint", web.HandleMemLeakRequest).
		WithPostHandler("/api/mem-footprint", web.HandleMemLeakRequest).
		Build()

	stopServerAsync := func() {
		server.StopAsync()
	}

	log.Println("Registering signal listeners for graceful HTTP server shutdown..")
	sys.RegisterShutdownHook(sys.NewSignalHook(syscall.SIGTERM, stopServerAsync))
	sys.RegisterShutdownHook(sys.NewSignalHook(syscall.SIGKILL, stopServerAsync))

	return server
}

func docsHandler(c *gin.Context) {
	host := c.Request.Host
	items := struct {
		Title string
		Urls  []string
	}{
		Title: "Usage:",
		Urls: []string{
			fmt.Sprintf("http://%s/api/cpu-load?time-sec=30[&cores=2]", host),
			fmt.Sprintf("http://%s/api/mem-footprint?amount-mb=1000", host),
		},
	}
	c.HTML(200, "help.tpl", items)
}
