package http

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
}

type Server interface {
	StartAsync()
	StopAsync()
	StopNow(timeout time.Duration) error
}

type ServerBuilder interface {
	Build() Server
	WithGetHandler(path string, handler func(c *gin.Context)) ServerBuilder
	WithPostHandler(path string, handler func(c *gin.Context)) ServerBuilder
	WithHTMLTemplates(templates ...string) ServerBuilder
}

type serverBuilder struct {
	port   int
	engine *gin.Engine
}

func (sb *serverBuilder) Build() Server {
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(sb.port),
		Handler: sb.engine,
	}

	s := &server{
		stopChan:   make(chan bool, 1),
		httpServer: httpServer,
	}

	return s
}

func (sb *serverBuilder) WithGetHandler(path string, handler func(c *gin.Context)) ServerBuilder {
	sb.engine.GET(path, handlerWrapperFor(path, handler))
	return sb
}

func (sb *serverBuilder) WithPostHandler(path string, handler func(c *gin.Context)) ServerBuilder {
	sb.engine.POST(path, handlerWrapperFor(path, handler))
	return sb
}

func (sb *serverBuilder) WithHTMLTemplates(templates ...string) ServerBuilder {
	sb.engine.LoadHTMLFiles(templates...)
	return sb
}

func NewServer(port int) ServerBuilder {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.HandleMethodNotAllowed = true

	sb := &serverBuilder{
		port:   port,
		engine: router,
	}

	return sb
}

func handlerWrapperFor(path string, handler func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Printf("Handling request at %s", path)
		// calling actual handler
		handler(c)

	}
}
