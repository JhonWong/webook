package web

import (
	"fmt"
	"net/http"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/service"
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
}

func NewUserHandler(us service.UserService, cs service.CodeService) *UserHandler {
	return &UserHandler{
		svc:              us,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		codeSvc:          cs,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/logout", u.Logout)
	ug.POST("/edit", u.Edit)
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("login_sms/code/send", u.SendSMSLoginCode)
	ug.POST("login_sms", u.LoginSMS)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassWord string `json:"confirmPassWord"`
		PassWord        string `json:"passWord"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	if req.PassWord != req.ConfirmPassWord {
		ctx.String(http.StatusOK, "两次输入密码不一致")
		return
	}

	ok, err = u.passwordRegexExp.MatchString(req.ConfirmPassWord)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码格式错误")
		return
	}

	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		PassWord: req.PassWord,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱已存在")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
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
	u.setJWTToken(ctx, user.Id)

	ctx.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) Logout(ctx *gin.Context) {
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
	claim, ok := value.(*UserClaim)
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
	u.setJWTToken(ctx, user.Id)
	ctx.JSON(http.StatusOK, Result{Msg: "登录成功"})
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, id int64) {
	claims := UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserId:    id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
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

type UserClaim struct {
	jwt.RegisteredClaims
	UserId    int64
	UserAgent string
}
