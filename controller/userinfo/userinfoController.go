package userinfo

import (
	"github.com/gin-gonic/gin"
	"mock_douyin_project/dao"
	"mock_douyin_project/service/user"
	"net/http"
)

type UserResponse struct {
	dao.CommonResponse
	User *dao.UserInfo `json:"user"`
}

func GetUserInfoByToken(c *gin.Context) {

	// 从jwt中间件中得到userid信息
	rawId, ok1 := c.Get("user_id")
	userId, ok2 := rawId.(int64)
	if !ok1 || !ok2 {
		c.JSON(http.StatusOK, UserResponse{
			CommonResponse: dao.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "解析用户id出现错误",
			},
		})

	}
	var userinfo dao.UserInfo
	err := user.DoQueryUserInfoByUserId(userId, &userinfo)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			CommonResponse: dao.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "查询用户信息出现错误",
			},
		})
	}
	c.JSON(http.StatusOK, UserResponse{
		CommonResponse: dao.CommonResponse{
			StatusCode: 0,
		},
		User: &userinfo,
	})
}
