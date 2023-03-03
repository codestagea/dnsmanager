package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"github.com/codestagea/bindmgr/internal/model"
	"github.com/codestagea/bindmgr/internal/server/middleware"
	"github.com/codestagea/bindmgr/internal/store"
	"github.com/codestagea/bindmgr/internal/tools"
	"strconv"
)

type DnsZoneHandler struct {
	authMiddleware *middleware.JwtMiddleware
}

func (h *DnsZoneHandler) InitRoute(r *gin.RouterGroup) {
	subRoute := r.Group("/v1/zones").Use(h.authMiddleware.MiddlewareFunc)
	subRoute.GET("", h.ListZones)
	subRoute.GET("/:id", h.ZoneDetail)
	subRoute.POST("", h.AddZone)
	subRoute.POST("/:id", h.UpdateZone)
}

func (h *DnsZoneHandler) ListZones(c *gin.Context) {
	pageQuery := model.NewPageQuery(c)
	search := c.Request.FormValue("search")

	if data, total, err := store.DnsZones.ListPage(search, pageQuery); err != nil {
		tools.RetOfErrMsg(c, 400, err.Error())
	} else {
		tools.Ok(c, model.NewPaged(data, total, pageQuery))
	}
}

func (h *DnsZoneHandler) ZoneDetail(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		logrus.Errorf("parse zone id %s fail: %s", c.Param("id"), err.Error())
		tools.RetOfErrMsg(c, 400, "zone id invalid")
	} else {
		if dn, dbErr := store.DnsZones.GetById(id); dbErr != nil {
			if errors.Is(dbErr, gorm.ErrRecordNotFound) {
				tools.RetOfErrMsg(c, 404, "zone not exist")
			} else {
				tools.RetOfErrMsg(c, 400, "query zone fail due to "+dbErr.Error())
			}
		} else {
			tools.Ok(c, dn)
		}
	}
}

func (h *DnsZoneHandler) AddZone(c *gin.Context) {
	zone := store.DnsZone{}
	if err := c.ShouldBind(&zone); err != nil {
		logrus.Errorf("bind zone fail: %v", err)
		tools.RetOfErrMsg(c, 400, "parse zone fail: "+err.Error())
		return
	}

	if err := zone.Validate(); err != nil {
		tools.RetOfErrMsg(c, 400, err.Error())
		return
	}

	_, err := store.DnsZones.GetByName(zone.Zone)
	if err == nil {
		tools.RetOfErrMsg(c, 400, zone.Zone+" already exist")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tools.RetOfErrMsg(c, 400, "get zone fail")
		return
	}

	loginUser, _ := middleware.LoginUser(c)

	//添加
	if err := store.DnsZones.AddZone(&zone, loginUser.Username); err != nil {
		logrus.Errorf("add zone fail: %v", err)
		tools.RetOfErrMsg(c, 400, "add zone fail")
		return
	}
	tools.Ok(c, zone.ID)
}
func (h *DnsZoneHandler) UpdateZone(c *gin.Context) {
	zone := store.DnsZone{}
	if err := c.ShouldBind(&zone); err != nil {
		logrus.Errorf("bind zone fail: %v", err)
		tools.RetOfErrMsg(c, 400, "parse zone fail: "+err.Error())
		return
	}

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		logrus.Errorf("parse zone id %s fail: %s", c.Param("id"), err.Error())
		tools.RetOfErrMsg(c, 400, "zone id invalid")
	} else {
		zone.ID = id
	}

	if err := zone.Validate(); err != nil {
		tools.RetOfErrMsg(c, 400, err.Error())
		return
	}

	loginUser, _ := middleware.LoginUser(c)

	//添加
	if err := store.DnsZones.UpdateById(&zone, loginUser.Username); err != nil {
		tools.RetOfErrMsg(c, 400, "update zone fail")
		return
	}
	tools.Ok(c, zone.ID)
}
