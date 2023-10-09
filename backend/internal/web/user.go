package web

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/johnwongx/webook/backend/pkg/ginx"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"net/http"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/service"
	myjwt "github.com/johnwongx/webook/backend/internal/web/jwt"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

	bizLogin = "login"
)

type UserHandler struct {
	svc              service.UserService
	codeSvc          service.CodeService
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	logger           logger.Logger

	myjwt.JwtHandler
}

func NewUserHandler(us service.UserService, cs service.CodeService,
	logger logger.Logger, j myjwt.JwtHandler) *UserHandler {
	return &UserHandler{
		svc:              us,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		codeSvc:          cs,
		logger:           logger,
		JwtHandler:       j,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", ginx.WrapReq[signUpReq](u.SignUp, u.logger))
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/refresh_token", u.RefreshToken)
	ug.POST("/logout", u.Logout)
	ug.POST("/edit", u.Edit)
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("login_sms/code/send", u.SendSMSLoginCode)
	ug.POST("login_sms", u.LoginSMS)
}

func (u *UserHandler) SignUp(ctx *gin.Context, req signUpReq) (ginx.Result, error) {
	ok, err := u.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, errors.New("系统错误")
	}
	if !ok {
		return ginx.Result{
			Code: 5,
			Msg:  "邮箱格式错误",
		}, errors.New("邮箱格式错误")
	}

	if req.PassWord != req.ConfirmPassWord {
		return ginx.Result{
			Code: 5,
			Msg:  "两次输入密码不一致",
		}, errors.New("两次输入密码不一致")
	}

	ok, err = u.passwordRegexExp.MatchString(req.ConfirmPassWord)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, errors.New("系统错误")
	}
	if !ok {
		return ginx.Result{
			Code: 5,
			Msg:  "密码格式错误",
		}, errors.New("密码格式错误")
	}

	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		PassWord: req.PassWord,
	})
	if err == service.ErrUserDuplicateEmail {
		return ginx.Result{
			Code: 5,
			Msg:  "邮箱已存在",
		}, errors.New("邮箱已存在")
	}
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, errors.New("系统错误")
	}

	return ginx.Result{
		Code: 1,
		Msg:  "注册成功",
	}, nil
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		PassWord string `json:"passWord"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.PassWord)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//设置session
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge: 10,
	})
	sess.Save()

	ctx.String(http.StatusOK, "Login Sucess")
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		PassWord string `json:"passWord"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.PassWord)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//设置session
	err = u.SetLoginToken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	tokenStr, err := u.ExtraToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	claims := &myjwt.RefreshClaim{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return myjwt.RtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	err = u.CheckSession(ctx, claims.SsId)
	if err != nil {
		//redis有问题，或者session无效
		//如果redis已经崩溃，可以考虑不在校验session
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//设置新的access token
	err = u.SetAccessToken(ctx, claims.UserId, claims.SsId)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "退出登录失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登录成功",
	})
}

func (u *UserHandler) LogoutSession(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()

	ctx.String(http.StatusOK, "Logout Sucess")
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	//校验信息
	type EditReq struct {
		NickName         string `json:"nickName"`
		Birthday         string `json:"birthday"`
		SelfIntroduction string `json:"selfIntroduction"`
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if len(req.NickName) > 60 {
		ctx.String(http.StatusOK, "名称长度不能大于60！")
		return
	}

	if !isValidBirthday(req.Birthday) {
		ctx.String(http.StatusOK, "生日格式非法，应为yyyy-dd-mm格式！")
		return
	}

	if len(req.SelfIntroduction) > 500 {
		ctx.String(http.StatusOK, "自我介绍长度不能超过500！")
		return
	}

	//读取用户id
	id, ok := getUserId(ctx)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	err := u.svc.Edit(ctx, id, req.NickName, req.Birthday, req.SelfIntroduction)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//保存信息
	ctx.String(http.StatusOK, "edit")
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	value, ok := ctx.Get("claims")
	claim, ok := value.(*myjwt.UserClaim)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	user, err := u.svc.Profile(ctx, claim.UserId)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	id, ok := getUserId(ctx)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	user, err := u.svc.Profile(ctx, id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (u *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type SendSMSReq struct {
		Phone string `json:"phoneNumber"`
	}

	var req SendSMSReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "请输入手机号码"})
		return
	}

	err := u.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{Msg: "发送成功"})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{Msg: "短信发送太频繁，请稍后再试"})
	default:
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type SMSReq struct {
		Phone string `json:"phoneNumber"`
		Code  string `json:"code"`
	}

	var req SMSReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ok, err := u.codeSvc.Verify(ctx, bizLogin, req.Code, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "验证码错误"})
		return
	}

	//验证码正确
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "系统错误"})
		return
	}
	u.SetLoginToken(ctx, user.Id)
	ctx.JSON(http.StatusOK, Result{Msg: "登录成功"})
}

func isValidBirthday(birthday string) bool {
	_, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
	}
	return err == nil
}

func getUserId(ctx *gin.Context) (int64, bool) {
	sess := sessions.Default(ctx)
	val := sess.Get("userId")
	id, ok := val.(int64)
	return id, ok
}

type signUpReq struct {
	Email           string `json:"email"`
	ConfirmPassWord string `json:"confirmPassWord"`
	PassWord        string `json:"passWord"`
}
