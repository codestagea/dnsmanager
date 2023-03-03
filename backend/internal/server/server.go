package server

import (
	"context"
	"github.com/codestagea/bindmgr/internal/metrics"
	"github.com/codestagea/bindmgr/internal/server/handler"
	"github.com/codestagea/bindmgr/internal/server/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/codestagea/bindmgr/global"
)

type DnsManagerServer struct {
	srv         *http.Server
	stopWatcher chan struct{}
}

func NewServer() (*DnsManagerServer, error) {
	httpConf := global.GVA_CONF.HttpConfig

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	middleware.InitMiddleware(r)
	authMiddleware := middleware.InitAuth(global.GVA_CONF.Jwt)
	handler.Init(r.Group(httpConf.ContextPath), authMiddleware)

	r.GET("/metrics", gin.WrapH(metrics.NewMetrics("dns-manager")))
	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, `{"status", "ok"}`)
	})

	srv := &http.Server{
		Addr:    httpConf.Host + ":" + strconv.Itoa(httpConf.Port),
		Handler: r,
	}
	return &DnsManagerServer{
		srv: srv,
	}, nil
}

func (cs *DnsManagerServer) Start() {
	logrus.Infof("start server at %s", cs.srv.Addr)

	cs.stopWatcher = make(chan struct{})

	go func() {
		// service connections
		if err := cs.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("listen: %s\n", err)
		}
	}()
}

func (ws *DnsManagerServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	close(ws.stopWatcher)
	if err := ws.srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server Shutdown: %v", err)
	}
}
