package relation

import (
	"github.com/gin-gonic/gin"
	"mock_douyin_project/service/relation"
	"net/http"
	"strconv"
)

func Action(c *gin.Context) {

	toUserIdStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type")

	userIdRaw, ok1 := c.Get("user_id")
	userId, ok2 := userIdRaw.(int64)
	toUserId, err1 := strconv.ParseInt(toUserIdStr, 10, 64)
	actionType, err2 := strconv.ParseInt(actionTypeStr, 10, 32)
	if !ok1 || !ok2 || err1 != nil || err2 != nil {
		c.JSON(400,
			relation.ActionResponse{StatusCode: -1, StatusMsg: "参数不匹配"})
		return
	}
	if userId == toUserId {
		c.JSON(400,
			relation.ActionResponse{StatusCode: -1, StatusMsg: "不能关注自己"})
		return
	}

	actionResponse := relation.Action(userId, toUserId, int32(actionType))
	if actionResponse.StatusCode == -1 {
		c.JSON(400, actionResponse)
		return
	}
	c.JSON(http.StatusOK, actionResponse)
}

func FollowList(c *gin.Context) {

	// 从jwt中间件中得到userid信息
	rawId, ok1 := c.Get("user_id")
	userId, ok2 := rawId.(int64)
	if !ok1 || !ok2 {
		c.JSON(400,
			relation.ActionResponse{StatusCode: -1, StatusMsg: "参数不匹配"})
		return
	}

	followListResponse := relation.FollowList(userId)
	if followListResponse.StatusCode == -1 {
		c.JSON(400, followListResponse)
		return
	}
	c.JSON(http.StatusOK, followListResponse)
}

func FollowerList(c *gin.Context) {

	// 从jwt中间件中得到userid信息
	rawId, ok1 := c.Get("user_id")
	userId, ok2 := rawId.(int64)
	if !ok1 || !ok2 {
		c.JSON(400,
			relation.ActionResponse{StatusCode: -1, StatusMsg: "参数不匹配"})
		return
	}

	followerListResponse := relation.FollowerList(userId)
	if followerListResponse.StatusCode == -1 {
		c.JSON(400, followerListResponse)
		return
	}
	c.JSON(http.StatusOK, followerListResponse)
}

func FriendList(c *gin.Context) {

	// 从jwt中间件中得到userid信息
	rawId, ok1 := c.Get("user_id")
	userId, ok2 := rawId.(int64)
	if !ok1 || !ok2 {
		c.JSON(400,
			relation.ActionResponse{StatusCode: -1, StatusMsg: "参数不匹配"})
		return
	}

	friendListResponse := relation.FriendList(userId)
	if friendListResponse.StatusCode == -1 {
		c.JSON(400, friendListResponse)
		return
	}
	c.JSON(http.StatusOK, friendListResponse)
}
