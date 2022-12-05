package dto

type AdminReq struct {
	Name     *string `json:"name" form:"name"`
	Username *string `json:"username" form:"username"`
	Email    *string `json:"email" form:"email"`
}

type AdminRequiredReq struct {
	Name     *string `json:"name" form:"name" binding:"required"`
	Username *string `json:"username" form:"username" binding:"required"`
	Email    *string `json:"email" form:"email" binding:"required"`
}

type AdminCreateReq struct {
	AdminRequiredReq
	Password *string `json:"password" form:"password"`
}

type AdminEditReq struct {
	ID       *uint64 `json:"id" form:"id" binding:"required"`
	Password *string `json:"password" form:"password"`
	AdminReq
}
