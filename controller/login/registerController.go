package login

import (
	"github.com/gin-gonic/gin"
	"mock_douyin_project/service/user"
)

func UserRegister(c *gin.Context) {

	user.RegisterService(c)
}
