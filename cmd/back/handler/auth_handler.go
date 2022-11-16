package handler

import (
	"cryptoshare/dto"
	"cryptoshare/middleware"
	"cryptoshare/repository"
	"cryptoshare/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	R    *gin.Engine
	repo *repository.Repository
}

func newAuthHandler(h *Handler) *authHandler {
	return &authHandler{
		R:    h.R,
		repo: h.repo,
	}
}

func (ctr *authHandler) register() {
	group := ctr.R.Group("/api/auth")
	group.POST("/login", ctr.login)
	group.Use(middleware.AuthMiddleware(ctr.repo))

	group.POST("/refresh", ctr.refresh)
	group.POST("/logout", ctr.logout)
}

func (ctr *authHandler) login(c *gin.Context) {
	req := dto.LoginReq{}
	res := &dto.Response{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	admin, err := ctr.repo.Admin.FindByField("email", req.Email)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	// otp validation
	if admin.OTPEnabled && req.OTP == "" {
		res.ErrCode = 400
		res.ErrMsg = "OTP is required."
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if admin.OTPEnabled {
		validOTP := utils.Validate2fa(req.OTP, admin.OTPSecret)
		if !validOTP {
			res.ErrCode = 400
			res.ErrMsg = "invalid otp"
			c.JSON(http.StatusBadRequest, res)
			return
		}
	}

	validPassword := utils.CheckPasswordHash(req.Password, admin.Password)
	if !validPassword {
		res.ErrCode = 400
		res.ErrMsg = "invalid password"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	accessToken, err := utils.GenerateAccessToken(admin.Username, true)
	if err != nil {
		res := utils.GenerateServerError(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	c.SetCookie("token", accessToken, int(time.Minute)*5, "/", c.Request.Host, true, true)

	res = utils.GenerateSuccessResponse(accessToken)
	c.JSON(res.HttpStatusCode, res)

}

func (ctr *authHandler) refresh(c *gin.Context) {
	token, _ := c.Cookie("token")
	refreshToken, err := utils.GenerateRefreshToken(token)
	if err != nil {
		res := utils.GenerateServerError(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	c.SetCookie("token", refreshToken, int(time.Minute)*5, "/", c.Request.Host, true, true)

	res := utils.GenerateSuccessResponse(refreshToken)
	c.JSON(res.HttpStatusCode, res)

}

func (ctr *authHandler) logout(c *gin.Context) {
	// immediately clear the token cookie
	c.SetCookie("token", "", 0, "/", c.Request.Host, true, true)

	res := utils.GenerateSuccessResponse("successfully logged out")
	c.JSON(res.HttpStatusCode, res)
}
