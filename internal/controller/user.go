package controller

import (
	"CampusTrader/internal/common/response"
	"CampusTrader/internal/model"
	"CampusTrader/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginReq struct {
	Account    string `json:"account" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

type RegisterReq struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required,min=6"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone" binding:"omitempty,numeric"`
	Nickname   string `json:"nick_name"`
}

type MeResp struct {
	model.User
}

type UserController struct {
	svc *service.UserService
}

func NewUserController(svc *service.UserService) *UserController {
	return &UserController{svc: svc}
}

func (c *UserController) Register(ctx *gin.Context) {
	var registerReq RegisterReq
	if err := ctx.ShouldBind(&registerReq); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if registerReq.Nickname == "" {
		registerReq.Nickname = registerReq.Username
	}

	if err := c.svc.Register(ctx, &model.User{
		Username: registerReq.Username,
		Password: registerReq.Password,
		Email:    registerReq.Email,
		Nickname: registerReq.Nickname,
	}); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *UserController) Login(ctx *gin.Context) {
	var loginReq LoginReq
	if err := ctx.ShouldBind(&loginReq); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	token, err := c.svc.Login(ctx, loginReq.Account, loginReq.Password)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(ctx, map[string]string{
		"token": token,
	})
}

// Me GET /api/v1/user/me need login
func (c *UserController) Me(ctx *gin.Context) {
	id := ctx.GetUint("userID")
	info, err := c.svc.GetUserInfo(ctx, id)
	if err != nil {
		return
	}
	response.Success(ctx, MeResp{*info})
}
