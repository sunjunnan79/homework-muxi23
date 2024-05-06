package request

type CheckPostMessageReq struct {
	Id string `json:"id" form:"id" binding:"required"`
}

type CheckLikeMessageReq struct {
	Id string `json:"id" form:"id" binding:"required"`
}
