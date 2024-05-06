package request

type CreatePostReq struct {
	TUuid string `json:"tUuid" form:"tUuid" binding:"required"`
	Body  string `json:"body" form:"body" binding:"required"`
}

type FindPostReq struct {
	PUuid string `json:"pUuid" form:"pUuid" binding:"required"`
}

type GetPostsReq struct {
	TUuid string `json:"tUuid" form:"tUuid" binding:"required"`
}
