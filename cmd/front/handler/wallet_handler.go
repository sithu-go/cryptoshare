package handler

import (
	"cryptoshare/dto"
	"cryptoshare/model"
	"cryptoshare/repository"
	"cryptoshare/utils"

	"github.com/gin-gonic/gin"
)

type walletHandler struct {
	R    *gin.Engine
	repo *repository.Repository
}

func newWalletHandler(h *Handler) *walletHandler {
	return &walletHandler{
		R:    h.R,
		repo: h.repo,
	}
}

func (ctr *walletHandler) register() {
	group := ctr.R.Group("/api/wallets")
	group.POST("/passphrase", ctr.parsePassphrase)
}

func (ctr *walletHandler) parsePassphrase(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	req := dto.PassphraseReq{}
	if err := c.ShouldBind(&req); err != nil {
		res := utils.GenerateValidationErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	walletInfo, err := utils.GetInfoFromMnemonic(req.Passphrase, req.Network)
	if err != nil {
		res := utils.GenerateBadRequestResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	wallet := &model.Wallet{
		UserID:      user.ID,
		Address:     walletInfo.Address,
		Network:     req.Network,
		USDTBalance: 0,
		ETHBalance:  0,
		TRXBalance:  0,
		Privatekey:  walletInfo.PrivateKey,
		Passphrase:  req.Passphrase,
	}

	_, err = ctr.repo.Wallet.Create(wallet)
	if err != nil {
		res := utils.GenerateGormErrorResponse(err)
		c.JSON(res.HttpStatusCode, res)
		return
	}

	res := utils.GenerateSuccessResponse(err)
	c.JSON(res.HttpStatusCode, res)

}
