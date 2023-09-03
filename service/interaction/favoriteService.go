package interactionService

import (
	"github.com/gin-gonic/gin"
	"log"
	"mock_douyin_project/dao"
	"net/http"
	"strconv"
)

// FavoriteAction 点赞以及取消点赞操作
// 通过对actionType字段进行switch判断产生分支
func FavoriteAction(c *gin.Context) {
	//先获取url中需要的参数并转化类型
	userIdAny, flag := c.Get("user_id")
	videoIdAny := c.Query("video_id")
	actionTypeAny := c.Query("action_type")
	userId, flag1 := userIdAny.(int64)
	videoId, flag2 := strconv.ParseInt(videoIdAny, 10, 64)
	actionType, flag3 := strconv.ParseInt(actionTypeAny, 10, 64)
	if !flag1 || flag2 != nil || flag3 != nil || !flag {
		c.JSON(http.StatusOK, dao.CommonResponse{
			StatusCode: 1,
			StatusMsg:  "参数传输错误",
		})
		return
	}
	favoriteDao := dao.NewFavoriteDao()
	switch actionType {
	case 1:
		{
			err := favoriteDao.FavoriteAction(userId, videoId)
			if err != nil {
				c.JSON(http.StatusOK, dao.CommonResponse{
					StatusCode: 0,
					StatusMsg:  "点赞成功"})
				return
			} else {
				log.Println(err)
				c.JSON(http.StatusOK, dao.CommonResponse{
					StatusCode: 2,
					StatusMsg:  "点赞失败，请刷新重试",
				})
				return
			}
		}
	case 2:
		{
			err := favoriteDao.FavoriteActionCancel(userId, videoId)
			if err != nil {
				c.JSON(http.StatusOK, dao.CommonResponse{
					StatusCode: 0,
					StatusMsg:  "取消点赞成功"})
				return
			} else {
				log.Println(err)
				c.JSON(http.StatusOK, dao.CommonResponse{
					StatusCode: 3,
					StatusMsg:  "取消点赞失败，请刷新重试",
				})
				return
			}
		}
	}
}

// FavoriteList 获取喜欢列表
func FavoriteList(c *gin.Context) {
	userIdAny, err1 := c.Get("user_id")
	userId, err2 := userIdAny.(int64)
	if !err1 || !err2 {
		c.JSON(http.StatusOK, dao.FavoriteListResponse{
			StatusCode: 1,
			StatusMsg:  "用户操作解析失败",
		})
		return
	}
	videos, err := dao.NewFavoriteDao().QueryFavoriteList(userId)
	if err != nil {
		c.JSON(http.StatusOK, dao.FavoriteListResponse{
			StatusCode: 2,
			StatusMsg:  "获取用户喜欢列表失败",
		})
		return
	}
	// 填充信息(作者信息和自己的点赞情况)
	for i := range videos {
		videos[i].IsFavorite = dao.NewVideoDAO().GetVideoIsFavoriteByUserIdAndVideoId(userId, videos[i].Id)
	}
	c.JSON(http.StatusOK, dao.FavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "获取用户喜欢列表成功",
		VideoList:  videos,
	})
	return
}
