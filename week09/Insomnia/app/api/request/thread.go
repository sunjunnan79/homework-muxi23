package request

type CreateThreadReq struct {
	Topic  string   `json:"topic" form:"topic" binding:"required"`
	Title  string   `json:"title" form:"title" binding:"required"`
	Body   string   `json:"body" form:"body" binding:"required"`
	Images []string `json:"images" form:"images"`
}

type FindThreadReq struct {
	TUuid string `json:"tUuid" form:"tUuid" binding:"required"`
}

type GetThreadsReq struct {
	Topic string `json:"topic" form:"topic" binding:"required"`
}
