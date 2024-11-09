package response

type LoginResponse struct {
	Token string `json:"token"`
}

type GetMyDataResponse struct {
	MyPost  uint `json:"myPost"`
	GetPost uint `json:"GetPost"`
	Likes   uint `json:"likes"`
}
