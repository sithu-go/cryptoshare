package middleware

import (
	"cryptoshare/repository"
	"cryptoshare/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type authHeader struct {
	AccessToken string `header:"Authorization"`
}

func AuthMiddleware(r *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		// bind Authorization Header to h and check for validation errors
		if err := c.ShouldBindHeader(&h); err != nil {
			res := utils.GenerateAuthErrorResponse(err)
			c.JSON(res.HttpStatusCode, res)
			c.Abort()
			return
		}

		accessToken := strings.Split(h.AccessToken, "Bearer ")

		if len(accessToken) != 2 {
			res := utils.GenerateAuthErrorResponse(fmt.Errorf("permission denied"))
			c.JSON(res.HttpStatusCode, res)
			c.Abort()
			return
		}

		// validate access token here
		claim, err := utils.ValidateAccessToken(accessToken[1])
		if err != nil {
			res := utils.GenerateAuthErrorResponse(nil)
			c.JSON(res.HttpStatusCode, res)
			c.Abort()
			return
		}

		if claim.IsAdmin {
			admin, err := r.Admin.FindByField("username", claim.Username)
			if err != nil {
				res := utils.GenerateGormErrorResponse(err)
				c.JSON(res.HttpStatusCode, res)
				c.Abort()
				return
			}
			c.Set("admin", admin)
		} else {
			user, err := r.User.FindByField("username", claim.Username)
			if err != nil {
				res := utils.GenerateGormErrorResponse(err)
				c.JSON(res.HttpStatusCode, res)
				c.Abort()
				return
			}
			c.Set("user", user)
		}
		c.Set("token", accessToken[1])
		c.Next()
	}
}
