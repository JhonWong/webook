package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/johnwongx/webook/backend/internal/service"
	"github.com/johnwongx/webook/backend/internal/service/oauth2/wechat"
	myjwt "github.com/johnwongx/webook/backend/internal/web/jwt"
	uuid "github.com/lithammer/shortuuid/v4"
)

type OAuth2WechatHandler struct {
	svc      wechat.Service
	userSvc  service.UserService
	stateKey []byte
	cfg      WechatHandlerConfig

	myjwt.JwtHandler
}

type WechatHandlerConfig struct {
	Secure bool
}

func NewWechatHandler(svc wechat.Service, userSvc service.UserService,
	cfg WechatHandlerConfig) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:      svc,
		userSvc:  userSvc,
		stateKey: []byte("95osj3fUD7fo0mlYdDbxcXz4VD2igvf1"),
		cfg:      cfg,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "生成认证url失败",
		})
		return
	}

	if err = h.setStateCookie(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		state: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.SetCookie("jwt-state", tokenStr, 600, "/oauth2/wechat/callback", "",
		h.cfg.Secure, true)
	return nil
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	//取出code与state(access_token)
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})
		return
	}

	//通过token获取wechat相关信息
	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	//通过信息查找用户
	//用户不存在时创建用户
	user, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	err = h.SetLoginToken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	cookie, err := ctx.Request.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("state cookie not found, %w", err)
	}

	claims := &StateClaim{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil {
		return err
	}

	if token == nil || !token.Valid {
		return fmt.Errorf("token expiration")
	}
	state := ctx.Query("state")
	if claims.state != state {
		return fmt.Errorf("token not equal")
	}
	return nil
}

type StateClaim struct {
	jwt.RegisteredClaims
	state string
}
