package dto

type RequestPayload struct {
	Page     int `json:"page" form:"page" binding:"required"`
	PageSize int `json:"page_size" form:"page_size"`
}

type Response struct {
	ErrCode        uint64 `json:"err_code"`
	ErrMsg         string `json:"err_msg"`
	Data           any    `json:"data,omitempty"`
	HttpStatusCode int    `json:"-"`
}

type ReqByID struct {
	ID uint64 `json:"id" form:"id" binding:"required"`
}

type ReqByIDs struct {
	IDS []uint64 `json:"ids" form:"ids" binding:"required,gte=1"`
}

type OTPReq struct {
	OTP string `json:"otp" binding:"required"`
}
