package interaction

import (
	"github.com/gin-gonic/gin"
	"mock_douyin_project/service/interactionService"
)

func CommentAction(c *gin.Context) {
	interactionService.CommentAction(c)
}

func CommentList(c *gin.Context) {
	interactionService.CommentList(c)
}
