package handler

import (
	"cryptoshare/dto"
	"cryptoshare/middleware"
	"cryptoshare/repository"
	"cryptoshare/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type bankHandler struct {
	R    *gin.Engine
	repo *repository.Repository
}

func newBankHandler(h *Handler) *bankHandler {
	return &bankHandler{
		R:    h.R,
		repo: h.repo,
	}
}

func (ctr *bankHandler) register() {
	group := ctr.R.Group("/api/banks")
	group.Use(middleware.AuthMiddleware(ctr.repo))
	group.GET("", ctr.getBanks)
	group.POST("", ctr.addBank)
	group.PATCH("", ctr.editBank)
	group.DELETE("", ctr.deleteBanks)
}

func (ctr *bankHandler) getBanks(c *gin.Context) {
	req := dto.RequestPayload{}
	if err := c.ShouldBindQuery(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	banks, err := ctr.repo.Bank.FindAll(&req)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	for _, bank := range banks {
		bank.PrivateKey = utils.GenerateRepeatedLetter("*", 64)
	}

	res := &dto.Response{
		ErrCode: 0,
		ErrMsg:  "Success",
		Data: gin.H{
			"data": banks,
		},
	}
	c.JSON(http.StatusOK, res)

}

func (ctr *bankHandler) addBank(c *gin.Context) {
	req := dto.CreateBankReq{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	encrytedPrivateKey, err := utils.EncryptAES(req.PrivateKey)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	req.PrivateKey = encrytedPrivateKey

	bank := ctr.repo.Bank.TransformCreateBankModel(&req)

	_, err = ctr.repo.Bank.Create(bank)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	res := &dto.Response{
		ErrCode: 0,
		ErrMsg:  "Success",
	}
	c.JSON(http.StatusOK, res)
}

func (ctr *bankHandler) editBank(c *gin.Context) {
	req := dto.UpdateBankReq{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	_, err := ctr.repo.Bank.Update(&req)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	res := &dto.Response{
		ErrCode: 0,
		ErrMsg:  "Success",
	}
	c.JSON(http.StatusOK, res)
}

func (ctr *bankHandler) deleteBanks(c *gin.Context) {
	req := dto.ReqByIDs{}
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err, "SDDD")
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	fmt.Printf("%+v\n", req)

	ids := utils.IdsIntToInCon(req.IDS)

	err := ctr.repo.Bank.DeleteMany(ids)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	res := &dto.Response{
		ErrCode: 0,
		ErrMsg:  "Success",
	}
	c.JSON(http.StatusOK, res)

}
