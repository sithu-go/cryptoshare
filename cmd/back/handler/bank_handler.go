package handler

import (
	"cryptoshare/dto"
	"cryptoshare/middleware"
	"cryptoshare/model"
	"cryptoshare/repository"
	"cryptoshare/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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

	// group.Use(middleware.OTPMiddleware("admin"))
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
		*bank.PrivateKey = utils.GenerateRepeatedLetter("*", 64)
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

	bank := model.Bank{}
	if err := copier.Copy(&bank, &req); err != nil {
		res := utils.GenerateServerError(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	encrytedPrivateKey, err := utils.EncryptAES(*bank.PrivateKey)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	bank.PrivateKey = &encrytedPrivateKey

	scanRecord := ctr.repo.Bank.GetAddressScanRecord(*bank.AddressType, *bank.WalletAddress)
	bank.ScanRecord = &scanRecord

	if err := ctr.repo.Bank.Create(c.Request.Context(), &bank); err != nil {
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
	fmt.Println("DSSDS")
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("DSSDS<", err.Error())
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	log.Println("<>?????")
	bank := model.Bank{}
	log.Printf("%+v\n", req)

	log.Printf("%+v\n", bank)
	if err := copier.Copy(&bank, &req); err != nil {
		res := utils.GenerateServerError(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}
	if err := ctr.repo.Bank.Update(c.Request.Context(), &bank); err != nil {
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
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

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
