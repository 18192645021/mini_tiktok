package interaction

import (
	"github.com/gin-gonic/gin"
	"mock_douyin_project/service/interaction"
)

func FavoriteAction(c *gin.Context) {
	interactionService.FavoriteAction(c)
}

func FavoriteList(c *gin.Context) {
	interactionService.FavoriteList(c)
}
