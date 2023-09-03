package relation

import (
	"mock_douyin_project/dao"
)

type ActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FollowListResponse struct {
	StatusCode int32           `json:"status_code"`
	StatusMsg  string          `json:"status_msg"`
	UserList   *[]dao.UserInfo `json:"user_list"`
}

type FriendUser struct {
	*dao.UserInfo
	Message string `json:"message"`
	MsgType int64  `json:"msgType"`
}

type FriendListResponse struct {
	StatusCode int32         `json:"status_code"`
	StatusMsg  string        `json:"status_msg"`
	UserList   *[]FriendUser `json:"user_list"`
}

// Action 关注与取消关注
func Action(userId int64, toUserId int64, actionType int32) *ActionResponse {
	isExist1 := dao.NewUserInfoDAO().IsUserExistById(userId)
	isExist2 := dao.NewUserInfoDAO().IsUserExistById(toUserId)
	if !isExist1 || !isExist2 {
		return &ActionResponse{StatusCode: -1, StatusMsg: "用户不存在"}
	}
	if actionType != 1 && actionType != 2 {
		return &ActionResponse{StatusCode: -1, StatusMsg: "请求操作类型错误"}
	}
	relationIsExist := dao.NewUserRelationDao().IsUserRelationExistByTwoId(userId, toUserId)
	if actionType == 1 {
		if relationIsExist {
			return &ActionResponse{StatusCode: -1, StatusMsg: "您已关注，请勿重复关注"}
		}
		userRelation := dao.UserRelation{UserInfoId: userId, FollowId: toUserId}
		err3 := dao.NewUserRelationDao().AddUserRelation(&userRelation, userId, toUserId)
		if err3 != nil {
			return &ActionResponse{StatusCode: -1, StatusMsg: "服务器错误"}
		}
		return &ActionResponse{StatusCode: 0, StatusMsg: "关注成功"}
	} else {
		if !relationIsExist {
			return &ActionResponse{StatusCode: -1, StatusMsg: "您未关注，请先关注"}
		}
		userRelation := dao.UserRelation{UserInfoId: userId, FollowId: toUserId}
		err3 := dao.NewUserRelationDao().DeleteUserRelation(&userRelation, userId, toUserId)
		if err3 != nil {
			return &ActionResponse{StatusCode: -1, StatusMsg: "服务器错误"}
		}
		return &ActionResponse{StatusCode: 0, StatusMsg: "取消关注成功"}
	}
}

// FollowList 获取关注用户列表
func FollowList(userId int64) *FollowListResponse {
	isExist := dao.NewUserInfoDAO().IsUserExistById(userId)
	if !isExist {
		return &FollowListResponse{StatusCode: -1, StatusMsg: "用户不存在"}
	}
	var followRelationList []dao.UserRelation
	err := dao.NewUserRelationDao().GetFollowRelationList(userId, &followRelationList)
	if err != nil {
		return &FollowListResponse{StatusCode: -1, StatusMsg: "服务器错误"}
	}
	var followUserList []dao.UserInfo
	for i := 0; i < len(followRelationList); i++ {
		followId := followRelationList[i].FollowId
		var followUser dao.UserInfo
		err = dao.NewUserInfoDAO().QueryUserInfoById(followId, &followUser)
		if err != nil {
			return &FollowListResponse{StatusCode: -1, StatusMsg: "服务器错误"}
		}
		followUser.IsFollow = true
		followUserList = append(followUserList, followUser)
	}
	return &FollowListResponse{StatusCode: 0, StatusMsg: "请求成功", UserList: &followUserList}
}

// FollowerList 获取粉丝列表
func FollowerList(userId int64) *FollowListResponse {
	isExist := dao.NewUserInfoDAO().IsUserExistById(userId)
	if !isExist {
		return &FollowListResponse{StatusCode: -1, StatusMsg: "用户不存在"}
	}
	var followerRelationList []dao.UserRelation
	err := dao.NewUserRelationDao().GetFollowerRelationList(userId, &followerRelationList)
	if err != nil {
		return &FollowListResponse{StatusCode: -1, StatusMsg: "服务器错误"}
	}
	var followerUserList []dao.UserInfo
	for i := 0; i < len(followerRelationList); i++ {
		followerId := followerRelationList[i].UserInfoId
		var followUser dao.UserInfo
		err = dao.NewUserInfoDAO().QueryUserInfoById(followerId, &followUser)
		if err != nil {
			return &FollowListResponse{StatusCode: -1, StatusMsg: "服务器错误"}
		}
		relationIsExist := dao.NewUserRelationDao().IsUserRelationExistByTwoId(userId, followerId)
		followUser.IsFollow = relationIsExist
		followerUserList = append(followerUserList, followUser)
	}
	return &FollowListResponse{StatusCode: 0, StatusMsg: "请求成功", UserList: &followerUserList}
}

// FriendList 获取好友列表
func FriendList(userId int64) *FriendListResponse {
	isExist := dao.NewUserInfoDAO().IsUserExistById(userId)
	if !isExist {
		return &FriendListResponse{StatusCode: -1, StatusMsg: "用户不存在"}
	}
	var followRelationList []dao.UserRelation
	err := dao.NewUserRelationDao().GetFollowRelationList(userId, &followRelationList)
	if err != nil {
		return &FriendListResponse{StatusCode: -1, StatusMsg: "服务器错误"}
	}
	var friendUserList []FriendUser
	for i := 0; i < len(followRelationList); i++ {
		followId := followRelationList[i].FollowId
		if !dao.NewUserRelationDao().IsUserRelationExistByTwoId(followId, userId) {
			continue
		}
		var friendInfo dao.UserInfo
		err = dao.NewUserInfoDAO().QueryUserInfoById(followId, &friendInfo)
		if err != nil {
			return &FriendListResponse{StatusCode: -1, StatusMsg: "服务器错误"}
		}
		friendInfo.IsFollow = true
		var lastMessage dao.Message
		err = dao.NewMessageDao().SelectLastMessage(userId, followId, &lastMessage)
		if err != nil {
			return &FriendListResponse{StatusCode: -1, StatusMsg: "服务器错误"}
		}
		var msgType int64 = 0
		if lastMessage.Id == userId {
			msgType = 1
		}
		friendUserList = append(friendUserList, FriendUser{UserInfo: &friendInfo, Message: lastMessage.Content, MsgType: msgType})
	}
	return &FriendListResponse{StatusCode: 0, StatusMsg: "请求成功", UserList: &friendUserList}
}
