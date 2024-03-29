package handler

import (
	"cryptoshare/dto"
	"cryptoshare/middleware"
	"cryptoshare/model"
	"cryptoshare/repository"
	"cryptoshare/utils"
	"log"
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
	group.POST("/register", ctr.singup)
	group.Use(middleware.AuthMiddleware(ctr.repo))

	group.POST("/refresh", ctr.refresh)
	group.POST("/logout", ctr.logout)
	group.POST("/generate/secret-key", ctr.generateSecretKey)
	group.POST("/enable/2fa", middleware.OTPMiddleware("admin"), ctr.enable2FactorAuth)
}
func (ctr *authHandler) singup(c *gin.Context) {
	req := dto.SingupReq{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	area, err := utils.GetArea(c.ClientIP())
	if err != nil {
		log.Println(err)
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		res := utils.GenerateServerError(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	user := &model.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: hash,
		IP:       c.ClientIP(),
		Location: area,
	}

	_, err = ctr.repo.User.Create(user)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)

}
func (ctr *authHandler) login(c *gin.Context) {
	req := dto.LoginReq{}
	res := &dto.Response{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	user, err := ctr.repo.User.FindByField("email", req.Email)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	// otp validation
	if user.OTPEnabled && req.OTP == "" {
		res.ErrCode = 400
		res.ErrMsg = "OTP is required."
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if user.OTPEnabled {
		validOTP := utils.Validate2fa(req.OTP, user.OTPSecret)
		if !validOTP {
			res.ErrCode = 400
			res.ErrMsg = "invalid otp"
			c.JSON(http.StatusBadRequest, res)
			return
		}
	}

	validPassword := utils.CheckPasswordHash(req.Password, user.Password)
	if !validPassword {
		res.ErrCode = 400
		res.ErrMsg = "invalid password"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.Username, false)
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

func (ctr *authHandler) generateSecretKey(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	key, err := utils.Create2fa(user.Username)
	if err != nil {
		res := utils.GenerateServerError(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	updateFields := &model.UpdateFields{
		Field: "id",
		Value: user.ID,
		Data: map[string]any{
			"otp_secret": key,
		},
	}
	_, err = ctr.repo.User.UpdateByFields(updateFields)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	res := utils.GenerateSuccessResponse(key)
	c.JSON(res.HttpStatusCode, res)
}

func (ctr *authHandler) enable2FactorAuth(c *gin.Context) {
	admin := c.MustGet("admin").(*model.Admin)

	updateFields := &model.UpdateFields{
		Field: "id",
		Value: admin.ID,
		Data: map[string]any{
			"otp_enabled": true,
		},
	}
	_, err := ctr.repo.Admin.UpdateByFields(c.Request.Context(), updateFields)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)
}
