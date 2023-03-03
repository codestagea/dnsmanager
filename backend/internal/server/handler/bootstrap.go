package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/codestagea/bindmgr/internal/server/middleware"
	// "github.com/codestagea/bindmgr/sync/tools"
)

type Handler interface {
	InitRoute(*gin.RouterGroup)
}

func Init(r *gin.RouterGroup, authMiddleware *middleware.JwtMiddleware) {
	handlers := []Handler{
		&UserHandler{authMiddleware: authMiddleware},
		&DnsZoneHandler{authMiddleware: authMiddleware},
		&DnsRecordHandler{authMiddleware: authMiddleware},
		&ViewHandler{authMiddleware: authMiddleware},
		&OperationLogHandler{authMiddleware: authMiddleware},
	}

	for _, h := range handlers {
		h.InitRoute(r)
	}
}
