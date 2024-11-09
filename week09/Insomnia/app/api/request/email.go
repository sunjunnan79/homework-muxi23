package request

type SendEmailRequest struct {
	Email string `json:"email" form:"email" binding:"required"`
}
