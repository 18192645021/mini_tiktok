package message

import (
	"mock_douyin_project/dao"
	"time"
)

type ActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type ChatResponse struct {
	StatusCode  int32          `json:"status_code"`
	StatusMsg   string         `json:"status_msg"`
	MessageList *[]dao.Message `json:"message_list"`
}

// Action 发送消息
func Action(userId int64, toUserId int64, actionType int32, content string) *ActionResponse {
	isExist := dao.NewUserInfoDAO().IsUserExistById(toUserId)
	if !isExist {
		return &ActionResponse{StatusCode: -1, StatusMsg: "目标用户不存在"}
	}
	if actionType != 1 {
		return &ActionResponse{StatusCode: -1, StatusMsg: "请求操作类型错误"}
	}

	followIsExist := dao.NewUserRelationDao().IsUserRelationExistByTwoId(userId, toUserId)
	followerIsExist := dao.NewUserRelationDao().IsUserRelationExistByTwoId(toUserId, userId)
	if !followIsExist || !followerIsExist {
		return &ActionResponse{StatusCode: -1, StatusMsg: "您和对方并非好友，不能发送消息"}
	}
	newMessage := dao.Message{FromUserId: userId, ToUserId: toUserId, Content: content, CreateTime: time.Now().Unix()}
	err := dao.NewMessageDao().AddMessage(&newMessage)
	if err != nil {
		return &ActionResponse{StatusCode: -1, StatusMsg: "服务器错误"}
	}
	return &ActionResponse{StatusCode: 0, StatusMsg: "消息发送成功"}
}

// ChatList 获取消息记录
func ChatList(userId int64, toUserId int64, preMsgTime int64) *ChatResponse {
	isExist := dao.NewUserInfoDAO().IsUserExistById(toUserId)
	if !isExist {
		return &ChatResponse{StatusCode: -1, StatusMsg: "目标用户不存在"}
	}

	var messageList []dao.Message
	err := dao.NewMessageDao().SelectMessageList(userId, toUserId, preMsgTime, &messageList)
	if err != nil {
		return &ChatResponse{StatusCode: -1, StatusMsg: "服务器错误"}
	}
	return &ChatResponse{StatusCode: 0, StatusMsg: "聊天记录访问成功", MessageList: &messageList}
}
