package dto

type CreateBankReq struct {
	Name          *string `json:"name" form:"name" binding:"required"`
	WalletAddress *string `json:"wallet_address" form:"wallet_address" binding:"required"`
	PrivateKey    *string `json:"private_key" form:"private_key" binding:"required"`
	AddressType   *string `json:"address_type" form:"address_type" binding:"required,oneof='ERC20' 'TRC20'"`
}

type UpdateBankReq struct {
	ReqByID
	Name          *string `json:"name" form:"name"`
	WalletAddress *string `json:"wallet_address" form:"wallet_address"`
	PrivateKey    *string `json:"private_key" form:"private_key"`
	AddressType   *string `json:"address_type" form:"address_type"` //  binding:"oneof=ERC20 TRC20"
}
