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

type DnsRecordHandler struct {
	authMiddleware *middleware.JwtMiddleware
}

func (h *DnsRecordHandler) InitRoute(r *gin.RouterGroup) {
	subRoute := r.Group("/v1/zone/:zoneId/records").Use(h.authMiddleware.MiddlewareFunc)
	subRoute.GET("", h.ListDnsRecords)
	subRoute.GET("/:id", h.DnsRecordDetail)
	subRoute.POST("", h.AddDnsRecord)
	subRoute.POST("/:id", h.UpdateDnsRecord)
}

func getZone(c *gin.Context) (*store.DnsZone, error) {
	if zoneId, err := strconv.ParseInt(c.Param("zoneId"), 10, 64); err != nil {
		logrus.Errorf("parse domain id %s fail: %s", c.Param("id"), err.Error())
		tools.RetOfErrMsg(c, 400, "domain id invalid")
		return nil, err
	} else {
		if d, dErr := store.DnsZones.GetById(zoneId); dErr != nil {
			if errors.Is(dErr, gorm.ErrRecordNotFound) {
				tools.RetOfErrMsg(c, 404, "domain not exist")
			} else {
				logrus.Errorf("get zone by id %d fail: %s", zoneId, dErr)
				tools.RetOfErrMsg(c, 400, err.Error())
			}
			return nil, dErr
		} else {
			return d, nil
		}
	}
}

func (h *DnsRecordHandler) ListDnsRecords(c *gin.Context) {
	zone, err := getZone(c)
	if err != nil {
		return
	}
	pageQuery := model.NewPageQuery(c)

	search := c.Request.FormValue("search")
	if data, total, err := store.DnsRecords.ListPage(zone.ID, search, pageQuery); err != nil {
		tools.RetOfErrMsg(c, 400, err.Error())
	} else {
		tools.Ok(c, model.NewPaged(data, total, pageQuery))
	}
}

func (h *DnsRecordHandler) DnsRecordDetail(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		logrus.Errorf("parse dns record id %s fail: %s", c.Param("id"), err.Error())
		tools.RetOfErrMsg(c, 400, "record id invalid")
	} else {
		if dn, dbErr := store.DnsRecords.GetById(id); dbErr != nil {
			if errors.Is(dbErr, gorm.ErrRecordNotFound) {
				tools.RetOfErrMsg(c, 404, "dns record not exist")
			} else {
				tools.RetOfErrMsg(c, 400, "query dns record fail due to "+dbErr.Error())
			}
		} else {
			tools.Ok(c, dn)
		}
	}
}

func (h *DnsRecordHandler) AddDnsRecord(c *gin.Context) {
	zone, err := getZone(c)
	if err != nil {
		return
	}
	record := store.DnsRecord{}
	if err := c.ShouldBind(&record); err != nil {
		logrus.Errorf("bind dns record fail: %v", err)
		tools.RetOfErrMsg(c, 400, "parse dns record fail: "+err.Error())
		return
	}
	record.ZoneId = zone.ID

	if err := record.Validate(); err != nil {
		tools.RetOfErrMsg(c, 400, err.Error())
		return
	}
	// 验证记录，是否为空，是否已经存在
	if record.Host == "" {
		tools.RetOfErrMsg(c, 400, "domain name should not empty")
		return
	}
	_, err = store.DnsRecords.GetByHost(zone.ID, record.Host, record.View)
	if err == nil {
		tools.RetOfErrMsg(c, 400, record.Host+" already exist")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tools.RetOfErrMsg(c, 400, "get domain record by host fail")
		return
	}

	loginUser, _ := middleware.LoginUser(c)

	//添加
	if err := store.DnsRecords.AddRecord(&record, loginUser.Username); err != nil {
		logrus.Errorf("add domain fail: %v", err)
		tools.RetOfErrMsg(c, 400, "add domain record fail")
		return
	}
	tools.Ok(c, record.ID)
}
func (h *DnsRecordHandler) UpdateDnsRecord(c *gin.Context) {
	zone, err := getZone(c)
	if err != nil {
		return
	}

	record := store.DnsRecord{}
	if err := c.ShouldBind(&record); err != nil {
		logrus.Errorf("bind domain name fail: %v", err)
		tools.RetOfErrMsg(c, 400, "parse domain name fail: "+err.Error())
		return
	}

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		logrus.Errorf("parse domain id %s fail: %s", c.Param("id"), err.Error())
		tools.RetOfErrMsg(c, 400, "domain id invalid")
	} else {
		record.ID = id
	}

	record.ZoneId = zone.ID

	if err := record.Validate(); err != nil {
		tools.RetOfErrMsg(c, 400, err.Error())
		return
	}
	//验证是否重名了
	if sameKey, err := store.DnsRecords.GetByHost(zone.ID, record.Host, record.View); err == nil {
		if sameKey.ID != record.ID { //改了host, view，导致重名了
			tools.RetOfErrMsg(c, 400, "same host and view already exist")
			return
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("get record by [zone: %d, host: %s, view: %s] fail: %v", zone.ID, record.Host, record.View, err.Error())
		tools.RetOfErrMsg(c, 400, "get record fail")
		return
	}

	loginUser, _ := middleware.LoginUser(c)

	//添加
	if err := store.DnsRecords.UpdateById(&record, loginUser.Username); err != nil {
		tools.RetOfErrMsg(c, 400, "update domain by name fail")
		return
	}
	tools.Ok(c, record.ID)
}
