package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/codestagea/bindmgr/internal/server/middleware"
	"github.com/codestagea/bindmgr/internal/store"
	"github.com/codestagea/bindmgr/internal/tools"
	"strconv"
)

type ViewHandler struct {
	authMiddleware *middleware.JwtMiddleware
}

func (h *ViewHandler) InitRoute(r *gin.RouterGroup) {
	subRoute := r.Group("/v1/view").Use(h.authMiddleware.MiddlewareFunc)
	subRoute.GET("", h.ListView)
	subRoute.POST("", h.AddView)
	subRoute.POST("/:id", h.UpdateView)
}

func (h *ViewHandler) ListView(c *gin.Context) {
	if data, err := store.Views.List(); err != nil {
		tools.RetOfErr(c, err)
	} else {
		tools.Ok(c, data)
	}
}

func (h *ViewHandler) AddView(c *gin.Context) {
	view := store.DnsView{}
	if err := c.ShouldBind(&view); err != nil {
		logrus.Errorf("bind view fail: %v", err)
		tools.RetOfErrMsg(c, 400, "parse view name fail: "+err.Error())
		return
	}
	if view.Name == "" {
		tools.RetOfErrMsg(c, 400, "view name cannot be null")
		return
	}

	loginUser, _ := middleware.LoginUser(c)

	//添加
	if err := store.Views.Add(&view, loginUser.Username); err != nil {
		tools.RetOfErrMsg(c, 400, err.Error())
		return
	}
	tools.Ok(c, view.ID)
}
func (h *ViewHandler) UpdateView(c *gin.Context) {
	view := store.DnsView{}
	if err := c.ShouldBind(&view); err != nil {
		logrus.Errorf("bind view fail: %v", err)
		tools.RetOfErrMsg(c, 400, "parse view fail: "+err.Error())
		return
	}

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		logrus.Errorf("parse zone id %s fail: %s", c.Param("id"), err.Error())
		tools.RetOfErrMsg(c, 400, "zone id invalid")
	} else {
		view.ID = id
	}

	loginUser, _ := middleware.LoginUser(c)

	//添加
	if err := store.Views.Update(&view, loginUser.Username); err != nil {
		tools.RetOfErrMsg(c, 400, "update zone by name fail")
		return
	}
	tools.Ok(c, view.ID)
}
