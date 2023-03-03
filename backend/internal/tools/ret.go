package tools

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/codestagea/bindmgr/internal/model"
	"net/http"
)

type Ret struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorWithCode struct {
	Code int
	Msg  string
}

var _ error = ErrorWithCode{400, "test"}

func (e ErrorWithCode) Error() string {
	return e.Msg
}

func Ok(c *gin.Context, data interface{}) {
	ret := Ret{
		Code:    0,
		Message: "",
		Data:    data,
	}
	c.JSON(http.StatusOK, ret)
}

func OkPaged(c *gin.Context, records interface{}, total int64, query *model.PageQuery) {
	paged := &model.Paged{
		Records:  records,
		Total:    total,
		PageNum:  query.PageNum,
		PageSize: query.PageSize,
	}
	Ok(c, paged)
}

func RetOfErr(c *gin.Context, err error) {
	RetOfErrMsg(c, retCode(err), err.Error())
}

func retCode(err error) int {
	if aErr, ok := err.(ErrorWithCode); ok {
		return aErr.Code
	}
	var ptrErr *ErrorWithCode
	if errors.As(err, &ptrErr) {
		return ptrErr.Code
	}
	return 400
}

func RetOfErrMsg(c *gin.Context, code int, message string) {
	ret := Ret{
		Code:    code,
		Message: message,
	}
	c.JSON(http.StatusOK, ret)
	c.Abort()
}
