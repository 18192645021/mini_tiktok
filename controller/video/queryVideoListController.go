package video

import (
	"github.com/gin-gonic/gin"
	"mock_douyin_project/dao"
	"mock_douyin_project/service/video"
	"net/http"
	"strconv"
)

type ListResponse struct {
	dao.CommonResponse
	*video.List
}

func QueryVideoList(c *gin.Context) {

	rawId1 := c.Query("user_id")
	userId, err := strconv.ParseInt(rawId1, 10, 64)
	if userId == 0 {
		rawId2, _ := c.Get("user_id")
		userId, _ = rawId2.(int64)
	}
	videoList, err := video.QueryVideoListByUserId(userId)
	if err != nil {
		QueryVideoListError(c, err.Error())
	}
	QueryVideoListSuccess(c, videoList)
}

func QueryVideoListError(c *gin.Context, message string) {
	c.JSON(http.StatusOK, ListResponse{
		CommonResponse: dao.CommonResponse{
			StatusCode: 1,
			StatusMsg:  message,
		},
	})
}
func QueryVideoListSuccess(c *gin.Context, videoList *video.List) {
	c.JSON(http.StatusOK, ListResponse{
		CommonResponse: dao.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "queryVideoListSuccess",
		},
		List: videoList,
	})
}
