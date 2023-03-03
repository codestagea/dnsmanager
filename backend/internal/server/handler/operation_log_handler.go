package handler

import (
	"github.com/codestagea/bindmgr/internal/model"
	"github.com/codestagea/bindmgr/internal/server/middleware"
	"github.com/codestagea/bindmgr/internal/store"
	"github.com/codestagea/bindmgr/internal/tools"
	"github.com/gin-gonic/gin"
	"strconv"
)

type OperationLogHandler struct {
	authMiddleware *middleware.JwtMiddleware
}

func (l *OperationLogHandler) InitRoute(r *gin.RouterGroup) {
	subRoute := r.Group("/v1/operation/logs")
	subRoute.GET("", l.ListOperationLog)
}

func (h *OperationLogHandler) ListOperationLog(c *gin.Context) {
	pageQuery := model.NewPageQuery(c)
	q := store.OperationLogQuery{}
	q.Search = c.Request.FormValue("search")
	q.TargetType = c.Request.FormValue("targetType")
	targetId := c.Request.FormValue("targetId")
	if targetId != "" {
		if tid, err := strconv.ParseInt(targetId, 10, 64); err != nil {
			q.TargetId = tid
		}
	}

	if data, total, err := store.OperationLogs.ListPage(&q, pageQuery); err != nil {
		tools.RetOfErrMsg(c, 400, err.Error())
	} else {
		tools.Ok(c, model.NewPaged(data, total, pageQuery))
	}
}
