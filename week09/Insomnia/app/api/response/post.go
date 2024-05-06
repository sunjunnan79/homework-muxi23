package response

type GetPostResponse struct {
	CreatedAt string `json:"createdAt"`
	PUuid     string `json:"pUuid"`
	TUuid     string `json:"tUuid"`
	UuId      string `json:"uuId"`
	Likes     uint   `json:"likes"`
	Body      string `json:"body"`
	Exist     string `json:"exist"`
}
