package response

type MessageResponse struct {
	CreatedAt string `json:"createdAt"`
	Id        string `json:"id"`
	PUuid     string `json:"pUuid"`
	TUuid     string `json:"tUuid"`
	UuId      string `json:"uuId"`
	RUuid     string `json:"rUuid"`
	Body      string `json:"body"`
	Title     string `json:"title"`
	Check     bool   `json:"check"`
}
