package response

type GetThreadResponse struct {
	CreatedAt string   `json:"createdAt"`
	TUuid     string   `json:"tUuid"`
	Topic     string   `json:"topic"`
	Title     string   `json:"title"`
	UuId      string   `json:"uuId"`
	Likes     uint     `json:"likes"`
	Body      string   `json:"body"`
	Number    uint     `json:"number"`
	Images    []string `json:"images"`
	Exist     string   `json:"exist"`
}
