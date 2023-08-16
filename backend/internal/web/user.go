package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/JhonWong/webook/backend/internal/domain"
	"github.com/JhonWong/webook/backend/internal/repository"
	"github.com/JhonWong/webook/backend/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	svc              *service.UserService
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:              svc,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/logout", u.Logout)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
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
	if err == repository.ErrUserDuplicateEmail {
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

	_, err := u.svc.Login(ctx, req.Email, req.PassWord)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//设置session
	token := jwt.New(jwt.SigningMethodHS512)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)

	ctx.String(http.StatusOK, "Login Sucess")
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
	id, err := getUserId(ctx)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	err = u.svc.Edit(ctx, id, req.NickName, req.Birthday, req.SelfIntroduction)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//保存信息
	ctx.String(http.StatusOK, "edit")
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
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

func isValidBirthday(birthday string) bool {
	_, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
	}
	return err == nil
}

func getUserId(ctx *gin.Context) (int64, error) {
	sess := sessions.Default(ctx)
	id := sess.Get("userId")
	switch v := id.(type) {
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to convert string to int64: %s", err)
		}
		return i, nil
	case int64:
		return id.(int64), nil
	default:
		return 0, fmt.Errorf("ussupported type: %T", id)
	}
}
