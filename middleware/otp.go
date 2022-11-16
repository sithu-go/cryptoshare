package middleware

import (
	"cryptoshare/dto"
	"cryptoshare/model"
	"cryptoshare/utils"

	"github.com/gin-gonic/gin"
)

func OTPMiddleware(userType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := dto.OTPReq{}
		if err := ctx.ShouldBind(&req); err != nil {
			res := utils.GenerateValidationErrorResponse(err)
			ctx.JSON(res.HttpStatusCode, res)
			return
		}
		// user or admin
		if userType == "admin" {
			admin := ctx.MustGet(userType).(*model.Admin)
			valid := utils.Validate2fa(req.OTP, admin.OTPSecret)
			if !valid {
				ctx.Abort()
				return
			}
			ctx.Next()
			return
		}
		// implement for user otp validation
		user := ctx.MustGet(userType).(*model.User)
		valid := utils.Validate2fa(req.OTP, user.OTPSecret)
		if !valid {
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
