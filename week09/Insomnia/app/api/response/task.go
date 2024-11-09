package response

type GetTaskResponse struct {
	Sum uint `json:"sum"`
	Day uint `json:"day"`
}

type GetAllTaskResponse struct {
	AllTask []GetTaskResponse `json:"allTask"`
}
