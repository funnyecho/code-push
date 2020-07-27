package http

import (
	"fmt"
	"github.com/funnyecho/code-push/gateway/client/interface/http/endpoints"
	"github.com/funnyecho/code-push/gateway/client/interface/http/middleware"
	"github.com/funnyecho/code-push/gateway/client/usecase"
	"github.com/gin-gonic/gin"
	stdHttp "net/http"
)

func New(uc usecase.UseCase, fns ...func(*Options)) *server {
	ctorOptions := &Options{}

	for _, fn := range fns {
		fn(ctorOptions)
	}

	svr := &server{
		uc:      uc,
		options: ctorOptions,
	}

	svr.initEndpoints()
	svr.initMiddleware()
	svr.initHttpHandler()

	return svr
}

type server struct {
	uc         usecase.UseCase
	options    *Options
	endpoints  *endpoints.Endpoints
	middleware *middleware.Middleware
	handler    stdHttp.Handler
}

func (s *server) ListenAndServe() error {
	addr := fmt.Sprintf(":%d", s.options.Port)
	server := &stdHttp.Server{
		Addr:           addr,
		Handler:        s,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}

func (s *server) ServeHTTP(writer stdHttp.ResponseWriter, request *stdHttp.Request) {
	s.handler.ServeHTTP(writer, request)
}

func (s *server) initEndpoints() {
	s.endpoints = endpoints.New(s.uc)
}

func (s *server) initMiddleware() {
	s.middleware = middleware.New(s.uc)
}

func (s *server) initHttpHandler() {
	r := gin.New()

	apiGroup := r.Group("/api")
	apiGroup.POST("/auth", s.endpoints.Auth)
	apiGroup.GET("/version/:version/download", s.endpoints.DownloadVersionPkg)
	apiGroup.GET("/version/:version/download/:filename", s.endpoints.DownloadVersionPkg)
	apiGroup.POST("/version/:version/upgrade", s.endpoints.VersionUpgradeQuery)

	s.handler = r
}

type Options struct {
	Port int
}