package request

type LikesReq struct {
	Uid string `json:"uid" form:"uid" binding:"required"`
}
