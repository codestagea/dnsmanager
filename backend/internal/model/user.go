package model

import "time"

type LoginUser struct {
	Uuid      string
	LoginType string

	Username    string
	UserId      int64
	ExpireTime  time.Time
	RefreshTime time.Time
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type UserInfo struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Perms    []string `json:"perms"`
}
