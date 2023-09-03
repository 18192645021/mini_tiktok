package interactionService

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mock_douyin_project/dao"
	"net/http"
	"strconv"
)

// CommentAction 进行添加评论和删除评论
// 通过对actionType字段判断进行分支
func CommentAction(c *gin.Context) {
	//先获取url中的参数并将转化类型
	userIdAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusOK, dao.CommentActionResponse{
			StatusCode: 2,
			StatusMsg:  "参数传输错误",
		})
		return
	}
	videoIdAny := c.Query("video_id")
	actionTypeAny := c.Query("action_type")
	userId, ok2 := userIdAny.(int64)
	videoId, err2 := strconv.ParseInt(videoIdAny, 10, 64)
	actionType, err3 := strconv.ParseInt(actionTypeAny, 10, 64)
	if !ok2 || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, dao.CommentActionResponse{
			StatusCode: 2,
			StatusMsg:  "参数传输错误",
		})
		return
	}
	commentDao := dao.NewCommentDAO()
	switch actionType {
	case 1:
		{
			commentText := c.Query("comment_text")
			if commentText == "" {
				c.JSON(http.StatusOK, dao.CommentActionResponse{
					StatusCode: 2,
					StatusMsg:  "评论不能为空",
				})
				return
			} else {
				comment := dao.Comment{UserInfoId: userId, VideoId: videoId, Content: commentText}
				if err := commentDao.AddCommentAndUpdateCount(&comment); err == nil {
					c.JSON(http.StatusOK, dao.CommentActionResponse{
						StatusCode: 0,
						StatusMsg:  "评论成功",
						Comment:    &comment,
					})
					return
				} else {
					c.JSON(http.StatusOK, dao.CommentActionResponse{
						StatusCode: 3,
						StatusMsg:  "评论失败，请刷新重试",
					})
				}
			}
		}
	case 2:
		{
			commentIdString := c.Query("comment_id")
			commentId, err5 := strconv.ParseInt(commentIdString, 10, 64)
			if err5 != nil {
				log.Println(err5)
				fmt.Print(commentId)
				c.JSON(http.StatusOK, dao.CommentActionResponse{
					StatusCode: 2,
					StatusMsg:  "评论ID参数传输错误",
				})
				return
			} else {
				if commentDao.IsCommentExistById(userId, videoId, commentId) {
					if err := commentDao.DeleteCommentAndUpdateCountById(commentId, videoId); err == nil {
						c.JSON(http.StatusOK, dao.CommentActionResponse{
							StatusCode: 0,
							StatusMsg:  "删除评论成功",
						})
						return
					} else {
						c.JSON(http.StatusOK, dao.CommentActionResponse{
							StatusCode: 5,
							StatusMsg:  "删除评论失败",
						})
						return
					}
				} else {
					c.JSON(http.StatusOK, dao.CommentActionResponse{
						StatusCode: 5,
						StatusMsg:  "评论或以被删除，请刷新重试",
					})
					return
				}
			}
		}

	}
}

// CommentList 获取某视频的评论列表
func CommentList(c *gin.Context) {
	videoIdString := c.Query("video_id")
	fmt.Println(videoIdString)
	videoId, err := strconv.ParseInt(videoIdString, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dao.CommentListResponse{
			StatusCode: 1,
			StatusMsg:  "参数传输错误",
		})
		return
	} else {
		if dao.NewVideoDAO().IsVideoExistById(videoId) {
			if commentList, err := dao.NewCommentDAO().QueryCommentListByVideoId(videoId); err != nil {
				c.JSON(http.StatusOK, dao.CommentListResponse{
					StatusCode: 2,
					StatusMsg:  "获取视频评论失败",
				})
				return
			} else {
				c.JSON(http.StatusOK, dao.CommentListResponse{
					StatusCode:  0,
					StatusMsg:   "获取视频评论成功",
					CommentList: commentList,
				})
				return
			}
		} else {
			c.JSON(http.StatusOK, dao.CommentListResponse{
				StatusCode: 2,
				StatusMsg:  "视频不存在，评论无法展示",
			})
		}
	}
}
