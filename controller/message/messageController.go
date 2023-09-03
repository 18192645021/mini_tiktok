package message

import (
	"github.com/gin-gonic/gin"
	"mock_douyin_project/service/message"
	"net/http"
	"strconv"
)

func Chat(c *gin.Context) {

	toUserIdStr := c.Query("to_user_id")
	preMsgTimeStr := c.Query("pre_msg_time")

	userIdRaw, ok1 := c.Get("user_id")
	userId, ok2 := userIdRaw.(int64)
	toUserId, err1 := strconv.ParseInt(toUserIdStr, 10, 64)
	preMsgTime, err2 := strconv.ParseInt(preMsgTimeStr, 10, 64)
	if !ok1 || !ok2 || err1 != nil || err2 != nil {
		c.JSON(400,
			message.ChatResponse{StatusCode: -1, StatusMsg: "参数不匹配"})
		return
	}

	chatResponse := message.ChatList(userId, toUserId, preMsgTime)
	if chatResponse.StatusCode == -1 {
		c.JSON(400, chatResponse)
		return
	}
	c.JSON(http.StatusOK, chatResponse)
}

func Action(c *gin.Context) {

	toUserIdStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type")
	content := c.Query("content")

	userIdRaw, ok1 := c.Get("user_id")
	userId, ok2 := userIdRaw.(int64)
	toUserId, err1 := strconv.ParseInt(toUserIdStr, 10, 64)
	actionType, err2 := strconv.ParseInt(actionTypeStr, 10, 32)
	if !ok1 || !ok2 || err1 != nil || err2 != nil {
		c.JSON(400,
			message.ActionResponse{StatusCode: -1, StatusMsg: "参数不匹配"})
		return
	}

	actionResponse := message.Action(userId, toUserId, int32(actionType), content)
	if actionResponse.StatusCode == -1 {
		c.JSON(400, actionResponse)
		return
	}
	c.JSON(http.StatusOK, actionResponse)
}
