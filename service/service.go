package service

import (
	"cryptoshare/conf"
)

var (
	tatumBaseURL = conf.TATUM_BASE_URL
	tatumApiKey  = conf.TATUM_API_KEY
)

type Service struct {
	TRC20 *trc20Service
	ERC20 *erc20Service
}

func NewService() *Service {
	trc20Service := newTRC20Service()
	erc20Service := newERC20Service()
	return &Service{
		TRC20: trc20Service,
		ERC20: erc20Service,
	}
}
