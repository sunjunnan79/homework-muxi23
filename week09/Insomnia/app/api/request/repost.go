package request

type CreateRePostReq struct {
	TUuid string `json:"tUuid" form:"tUuid" binding:"required"`
	PUuid string `json:"pUuid" form:"pUuid" binding:"required"`
	Body  string `json:"body" form:"body" binding:"required"`
}

type FindRePostReq struct {
	RUuid string `json:"rUuid" form:"rUuid" binding:"required"`
}

type GetRePostsReq struct {
	PUuid string `json:"pUuid" form:"pUuid" binding:"required"`
}
