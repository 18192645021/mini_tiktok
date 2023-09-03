package login

import (
	"github.com/gin-gonic/gin"
	"mock_douyin_project/service/user"
)

func UserLogin(c *gin.Context) {

	user.LoginService(c)
}
