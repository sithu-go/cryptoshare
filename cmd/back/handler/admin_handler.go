package handler

import (
	"cryptoshare/dto"
	"cryptoshare/middleware"
	"cryptoshare/model"
	"cryptoshare/repository"
	"cryptoshare/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type adminHandler struct {
	R    *gin.Engine
	repo *repository.Repository
}

func newAdminHandler(h *Handler) *adminHandler {
	return &adminHandler{
		R:    h.R,
		repo: h.repo,
	}
}

func (ctr *adminHandler) register() {

	group := ctr.R.Group("/api/admins")
	group.Use(middleware.AuthMiddleware(ctr.repo))

	group.GET("", ctr.getAdmins)
	group.POST("", ctr.createAdmin)
	group.PATCH("", ctr.updateAdmin)
	group.DELETE("", ctr.deleteAdmins)
	group.PUT("", ctr.recoverAdmin)
}

func (ctr *adminHandler) recoverAdmin(c *gin.Context) {
	req := dto.ReqByIDs{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	ids := utils.IdsIntToInCon(req.IDS)

	if err := ctr.repo.Admin.RecoverAdmins(c.Request.Context(), ids); err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	res := utils.GenerateSuccessResponse(nil)

	c.JSON(res.HttpStatusCode, res)
}

func (ctr *adminHandler) getAdmins(c *gin.Context) {
	req := dto.PageReq{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	list, total, err := ctr.repo.Admin.List(c.Request.Context(), &req)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	data := gin.H{
		"list":  list,
		"total": total,
	}
	res := utils.GenerateSuccessResponse(data)

	c.JSON(res.HttpStatusCode, res)
}

func (ctr *adminHandler) createAdmin(c *gin.Context) {
	req := dto.AdminCreateReq{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	admin := model.Admin{}
	if err := copier.Copy(&admin, &req); err != nil {
		res := utils.GenerateServerError(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	if err := ctr.repo.Admin.Create(c.Request.Context(), &admin); err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)

}

func (ctr *adminHandler) updateAdmin(c *gin.Context) {
	req := dto.AdminEditReq{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	admin := model.Admin{}
	if err := copier.Copy(&admin, &req); err != nil {
		res := utils.GenerateServerError(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	if admin.Password != nil {
		hashedPassword, err := utils.HashPassword(*admin.Password)
		log.Println(err)
		admin.Password = &hashedPassword
	}

	if err := ctr.repo.Admin.Update(c.Request.Context(), &admin); err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)
}

func (ctr *adminHandler) deleteAdmins(c *gin.Context) {
	req := dto.ReqByIDs{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	ids := utils.IdsIntToInCon(req.IDS)

	if err := ctr.repo.Admin.DeleteMany(c.Request.Context(), ids); err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	res := utils.GenerateSuccessResponse(nil)
	c.JSON(res.HttpStatusCode, res)
}
