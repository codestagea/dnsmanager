package middleware

import (
	"github.com/codestagea/bindmgr/config"
	"github.com/codestagea/bindmgr/internal/cache"
	"github.com/codestagea/bindmgr/internal/model"
	"github.com/codestagea/bindmgr/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	login_user_key  = "user_key"
	login_user_type = "user_type"

	token_prefix = "Bearer "
	token_header = "Authorization"

	sign_method = "HS512"

	LOGIN_USER_ATTR = "login_user"
)

var InvalidTokenError = tools.ErrorWithCode{3001, "invalid token"}
var TokenExpireError = tools.ErrorWithCode{3000, "token expire"}

func InitAuth(jwtCfg config.Jwt) *JwtMiddleware {
	return &JwtMiddleware{
		securityKey: []byte(jwtCfg.Secret),
		timeout:     jwtCfg.Timeout,
		cache:       cache.NewLocalTokenCache(1 * time.Second),
	}
}

type JwtMiddleware struct {
	securityKey []byte
	timeout     int
	cache       cache.TokenCache
}

func (jm *JwtMiddleware) MiddlewareFunc(c *gin.Context) {
	u, err := jm.GetUserInfo(c)

	if err != nil {
		logrus.Errorf("get user info fail: %s", err.Error())
		tools.RetOfErr(c, err)
		c.Abort()
		return
	}

	if time.Now().After(u.RefreshTime) {
		jm.refreshToken(u)
	}
	c.Set(LOGIN_USER_ATTR, u)
	c.Next()
}

func LoginUser(c *gin.Context) (*model.LoginUser, bool) {
	if v, exist := c.Get(LOGIN_USER_ATTR); !exist {
		return nil, false
	} else {
		return v.(*model.LoginUser), true
	}
}

func (jm *JwtMiddleware) GetUserInfo(c *gin.Context) (*model.LoginUser, error) {
	u := model.LoginUser{}
	token := GetToken(c)
	if token == "" {
		logrus.Errorf("token is emtpy")
		return nil, InvalidTokenError
	}
	if claims, err := jm.parseToken(token); err != nil {
		return nil, err
	} else {
		if u.Uuid, err = extractString(login_user_key, token, claims); err != nil {
			return nil, err
		}
		if u.LoginType, err = extractString(login_user_type, token, claims); err != nil {
			return nil, err
		}
		return jm.cache.GetLoginUser(u.Uuid)
	}
}

func (jm *JwtMiddleware) CreateToken(user *model.LoginUser) (string, error) {
	claims := jwt.MapClaims{}
	user.Uuid = uuid.NewString()

	if err := jm.refreshToken(user); err != nil {
		return "", err
	}

	claims[login_user_key] = user.Uuid
	claims[login_user_type] = user.LoginType

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	if ss, err := token.SignedString(jm.securityKey); err != nil {
		return "", errors.WithMessage(InternalError, "failed to create jwt: "+err.Error())
	} else {
		return token_prefix + ss, nil
	}
}

func (jm *JwtMiddleware) DelCache(user *model.LoginUser) error {
	return jm.cache.DelToken(user.Uuid)
}

func (jm *JwtMiddleware) parseToken(token string) (jwt.MapClaims, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, InvalidTokenError
		}

		return jm.securityKey, nil
	})
	if err != nil {
		if !errors.Is(err, InvalidTokenError) {
			logrus.Errorf("parse token fail: %v", err)
			return nil, InvalidTokenError
		}
		return nil, err
	}
	if !jwtToken.Valid {
		logrus.Errorf("token is invalid: %s", token)
		return nil, InvalidTokenError
	}

	claims, _ := jwtToken.Claims.(jwt.MapClaims)
	return claims, nil
}

func (jm *JwtMiddleware) refreshToken(u *model.LoginUser) error {
	d := time.Second * time.Duration(jm.timeout)
	now := time.Now()
	u.RefreshTime = now.Add(time.Duration(jm.timeout/2) * time.Second)
	u.ExpireTime = now.Add(d)
	return jm.cache.RefreshToken(u)
}

// GetToken help to get the JWT token string
func GetToken(c *gin.Context) string {
	token := c.GetHeader(token_header)
	if token == "" {
		return ""
	}

	return strings.TrimPrefix(token, token_prefix)
}

func extractString(key, token string, claims jwt.MapClaims) (string, error) {
	if claims[key] == nil {
		logrus.Errorf("%s key missing in token: %s", key, token)
		return "", InvalidTokenError
	}

	if v, ok := claims[key].(string); !ok {
		logrus.Errorf("%s key format error token: %s", key, token)
		return "", InvalidTokenError
	} else {
		return v, nil
	}
}
