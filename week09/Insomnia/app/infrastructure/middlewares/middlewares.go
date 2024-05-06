package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Load(r *gin.Engine) {
	UseRecovery(r)
}
