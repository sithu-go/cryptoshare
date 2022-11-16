package handler

import (
	"cryptoshare/dto"
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
}

func (ctr *authHandler) login(c *gin.Context) {
	req := dto.LoginReq{}
	res := dto.Response{}
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

	c.JSON(http.StatusOK, gin.H{
		"admin_token": accessToken,
	})

}
