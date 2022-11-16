package middleware

import (
	"cryptoshare/repository"
	"cryptoshare/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(r *repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		// res := &dto.Response{}
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				res := utils.GenerateAuthErrorResponse(err)
				ctx.JSON(res.HttpStatusCode, res)
				return
			}
			// For any other type of error, return a bad request status
			res := utils.GenerateBadRequestResponse(err)
			ctx.JSON(res.HttpStatusCode, res)
			return
		}

		claim, err := utils.ValidateAccessToken(token)
		if err != nil {
			res := utils.GenerateAuthErrorResponse(err)
			ctx.JSON(res.HttpStatusCode, res)
			return
		}
		if claim.IsAdmin {
			admin, err := r.Admin.FindByField("username", claim.Username)
			if err != nil {
				res := utils.GenerateGormErrorResponse(err)
				ctx.JSON(res.HttpStatusCode, res)
				return
			}
			ctx.Set("admin", admin)
		} else {
			// user, err := r.
			fmt.Println("will implement later")
		}
		ctx.Next()

	}
}
