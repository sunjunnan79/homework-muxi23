package controller

import (
	"Insomnia/app/api/response"
	"Insomnia/app/utility/tube"
	"github.com/gin-gonic/gin"
)

type Tube struct{}
type Token struct {
	Token string `json:"token"`
}

func (t *Tube) GetQNToken(c *gin.Context) {
	response.OkMsgData(c, "获取token成功", Token{Token: tube.GetQNToken()})
	return
}
