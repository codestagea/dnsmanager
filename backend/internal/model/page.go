package model

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Paged struct {
	Records  interface{} `json:"records"`
	Total    int64       `json:"total"`
	PageNum  int         `json:"pageIndex"`
	PageSize int         `json:"pageSize"`
}

func NewPaged(records interface{}, total int64, query *PageQuery) *Paged {
	return &Paged{
		Records:  records,
		Total:    total,
		PageNum:  query.PageNum,
		PageSize: query.PageSize,
	}
}

type PageQuery struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
	Offset   int `json:"offset"`
}

func NewPageQuery(c *gin.Context) *PageQuery {
	pageSize := 0
	pageSizeStr := c.Request.FormValue("pageSize")
	if pageSizeStr != "" {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	pageNum := 0
	pageNumStr := c.Request.FormValue("pageNum")
	if pageNumStr != "" {
		pageNum, _ = strconv.Atoi(pageNumStr)
	}

	if pageNum <= 0 {
		pageNum = 1
	}

	return &PageQuery{
		PageNum:  pageNum,
		PageSize: pageSize,
		Offset:   (pageNum - 1) * pageSize,
	}
}
