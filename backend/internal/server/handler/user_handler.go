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

type UserHandler struct {
	authMiddleware *middleware.JwtMiddleware
}

func (h *UserHandler) InitRoute(r *gin.RouterGroup) {
	subRoute := r.Group("/v1/users")
	subRoute.POST("auth/login", h.Login)
	subRoute.Use(h.authMiddleware.MiddlewareFunc).GET("me/info", h.CurrentUser)
	subRoute.Use(h.authMiddleware.MiddlewareFunc).POST("auth/logout", h.Logout)
	subRoute.Use(h.authMiddleware.MiddlewareFunc).POST(":username/status", h.ChangeUserStatus)
	subRoute.Use(h.authMiddleware.MiddlewareFunc).POST(":username/password/reset", h.ChangePasswordByAdmin)
	subRoute.Use(h.authMiddleware.MiddlewareFunc).POST("password/reset", h.ChangePassword)
	subRoute.Use(h.authMiddleware.MiddlewareFunc).POST("", h.AddUser)
	subRoute.Use(h.authMiddleware.MiddlewareFunc).POST(":username", h.UpdateUser)
	subRoute.Use(h.authMiddleware.MiddlewareFunc).GET("page", h.ListUserPage)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginReq

	if err := c.ShouldBind(&req); err != nil {
		tools.RetOfErr(c, err)
		return
	}

	user, err := store.Users.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tools.RetOfErrMsg(c, 404, "user not found")
		} else {
			logrus.Errorf("failed to query user %s: %v", req.Username, err)
			tools.RetOfErrMsg(c, 400, "get user fail")
		}
		return
	}
	if user.Status != 1 {
		tools.RetOfErrMsg(c, 401, "user is disabled")
		return
	}

	pwd, err := tools.PwdRsaDecode(req.Password)
	if err != nil {
		logrus.Errorf("decode %s password fail: %v", req.Username, err)
		tools.RetOfErrMsg(c, 400, "password invalid")
		return
	}
	if _, err = tools.HashPwdCompare(user.Password, pwd); err != nil {
		logrus.Errorf("%s password hash compare fail: %v", req.Username, err)
		tools.RetOfErrMsg(c, 400, "password invalid")
		return
	}

	token, err := h.authMiddleware.CreateToken(&model.LoginUser{
		Username:  req.Username,
		LoginType: "pwd",
		UserId:    1,
	})

	if err != nil {
		logrus.Errorf("failed to create token for user %s: %s", req.Username, err.Error())
		tools.RetOfErrMsg(c, 400, "login fail")
	} else {
		tools.Ok(c, map[string]string{"token": token})
	}
}

func (h *UserHandler) CurrentUser(c *gin.Context) {
	u, _ := middleware.LoginUser(c)

	info := model.UserInfo{
		Username: u.Username,
		Roles:    []string{"admin"},
		Perms:    []string{"admin"},
	}
	tools.Ok(c, info)
}

func (h *UserHandler) Logout(c *gin.Context) {
	u, _ := middleware.LoginUser(c)
	h.authMiddleware.DelCache(u)
	tools.Ok(c, "logout")
}

func (h *UserHandler) ListUserPage(c *gin.Context) {
	pageQuery := model.NewPageQuery(c)
	search := c.Request.FormValue("search")
	status := c.Request.FormValue("status")
	statusCode := -1
	if status != "" {
		statusCode, _ = strconv.Atoi(status)
	}

	q := store.UserQuery{
		Search: search,
		Status: statusCode,
	}
	if data, total, err := store.Users.ListUser(q, pageQuery); err != nil {
		tools.RetOfErr(c, err)
	} else {
		tools.OkPaged(c, data, total, pageQuery)
	}
}

func (h *UserHandler) ChangeUserStatus(c *gin.Context) {
	var req store.User
	if err := c.ShouldBind(&req); err != nil {
		tools.RetOfErr(c, err)
		return
	}
	username := c.Param("username")
	if _, err := store.Users.GetByUsername(username); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tools.RetOfErrMsg(c, 404, "user not found")
		} else {
			logrus.Errorf("failed to query user %s: %v", req.Username, err)
			tools.RetOfErrMsg(c, 400, "get user fail")
		}
	} else {
		if cErr := store.Users.ChangeStatus(username, req.Status); err != nil {
			logrus.Errorf("failed to change user status %s: %v", username, cErr)
			tools.RetOfErrMsg(c, 400, "change status fail")
		} else {
			tools.Ok(c, "ok")
		}
	}
}
func (h *UserHandler) ChangePasswordByAdmin(c *gin.Context) {
	var req store.User
	if err := c.ShouldBind(&req); err != nil {
		tools.RetOfErr(c, err)
		return
	}
	username := c.Param("username")
	if _, err := store.Users.GetByUsername(username); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tools.RetOfErrMsg(c, 404, "user not found")
		} else {
			logrus.Errorf("failed to query user %s: %v", req.Username, err)
			tools.RetOfErrMsg(c, 400, "get user fail")
		}
		return
	}

	if err := tools.VerifyPwdStrength(req.Password); err != nil {
		tools.RetOfErr(c, err)
		return
	}
	if pwd, err := tools.PwdHash(req.Password); err != nil {
		tools.RetOfErr(c, err)
	} else {
		if uErr := store.Users.ChangePwd(username, pwd); err != nil {
			tools.RetOfErr(c, uErr)
		} else {
			tools.Ok(c, "ok")
		}
	}
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req store.User
	if err := c.ShouldBind(&req); err != nil {
		tools.RetOfErr(c, err)
		return
	}
	loginUser, _ := middleware.LoginUser(c)
	if err := tools.VerifyPwdStrength(req.Password); err != nil {
		tools.RetOfErr(c, err)
		return
	}
	if pwd, err := tools.PwdHash(req.Password); err != nil {
		tools.RetOfErr(c, err)
	} else {
		if uErr := store.Users.ChangePwd(loginUser.Username, pwd); err != nil {
			tools.RetOfErr(c, uErr)
		} else {
			tools.Ok(c, "ok")
		}
	}
}

func (h *UserHandler) AddUser(c *gin.Context) {
	var req store.User
	if err := c.ShouldBind(&req); err != nil {
		tools.RetOfErr(c, err)
		return
	}

	if _, err := store.Users.GetByUsername(req.Username); err == nil {
		tools.RetOfErrMsg(c, 400, "username "+req.Username+" is already exist")
		return

	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("failed to query user %s: %v", req.Username, err)
		tools.RetOfErrMsg(c, 400, "get user fail")
		return
	}

	u := store.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Status:   req.Status,
	}

	if req.Password != "" {
		if err := tools.VerifyPwdStrength(req.Password); err != nil {
			tools.RetOfErr(c, err)
			return
		}
		if pwd, err := tools.PwdHash(req.Password); err != nil {
			tools.RetOfErr(c, err)
		} else {
			u.Password = pwd
		}
	}

	if uErr := store.Users.AddUser(&u); uErr != nil {
		tools.RetOfErr(c, uErr)
	} else {
		tools.Ok(c, "ok")
	}
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req store.User
	if err := c.ShouldBind(&req); err != nil {
		tools.RetOfErr(c, err)
		return
	}

	u := store.User{
		//Username: req.Username,
		Nickname: req.Nickname,
		Status:   req.Status,
	}

	username := c.Param("username")

	if exist, err := store.Users.GetByUsername(username); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tools.RetOfErrMsg(c, 404, "user not found")
		} else {
			logrus.Errorf("failed to query user %s: %v", username, err)
			tools.RetOfErrMsg(c, 400, "get user fail")
		}
		return
	} else {
		u.ID = exist.ID
	}

	if req.Password != "" {
		if err := tools.VerifyPwdStrength(req.Password); err != nil {
			tools.RetOfErr(c, err)
			return
		}
		if pwd, err := tools.PwdHash(req.Password); err != nil {
			tools.RetOfErr(c, err)
		} else {
			u.Password = pwd
		}
	}
	if uErr := store.Users.UpdateUser(&u); uErr != nil {
		tools.RetOfErr(c, uErr)
	} else {
		tools.Ok(c, "ok")
	}
}
