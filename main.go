package main

import (
	"github.com/gin-gonic/gin"
	"mock_douyin_project/controller/interaction"
	"mock_douyin_project/controller/login"
	"mock_douyin_project/controller/message"
	"mock_douyin_project/controller/relation"
	"mock_douyin_project/controller/userinfo"
	"mock_douyin_project/controller/video"
	"mock_douyin_project/dao"
	midddleware "mock_douyin_project/middleware"
)

func main() {

	r := gin.Default()
	dao.InitDB()
	r.Static("static", "./static")
	baseGroup := r.Group("/douyin")
	baseGroup.POST("/user/register/", midddleware.SHAMiddleWare(), login.UserRegister)
	baseGroup.GET("/user/", midddleware.JwtAuthMiddleware(), userinfo.GetUserInfoByToken)
	baseGroup.POST("/user/login/", midddleware.SHAMiddleWare(), login.UserLogin)
	baseGroup.GET("/publish/list/", midddleware.JwtAuthMiddleware(), video.QueryVideoList)
	baseGroup.GET("/feed/", video.FeedVideoList)
	baseGroup.POST("/publish/action/", midddleware.JwtAuthMiddleware(), video.PostVideo)

	//互动接口实现
	baseGroup.POST("/favorite/action/", midddleware.JwtAuthMiddleware(), interaction.FavoriteAction)
	baseGroup.GET("/favorite/list/", midddleware.JwtAuthMiddleware(), interaction.FavoriteList)
	baseGroup.POST("/comment/action/", midddleware.JwtAuthMiddleware(), interaction.CommentAction)
	baseGroup.GET("/comment/list/", interaction.CommentList)

	// 社交
	baseGroup.POST("/relation/action/", midddleware.JwtAuthMiddleware(), relation.Action)
	baseGroup.GET("/relation/follow/list/", midddleware.JwtAuthMiddleware(), relation.FollowList)
	baseGroup.GET("/relation/follower/list/", midddleware.JwtAuthMiddleware(), relation.FollowerList)
	baseGroup.GET("/relation/friend/list/", midddleware.JwtAuthMiddleware(), relation.FriendList)
	baseGroup.GET("/message/chat/", midddleware.JwtAuthMiddleware(), message.Chat)
	baseGroup.POST("/message/action/", midddleware.JwtAuthMiddleware(), message.Action)
	err := r.Run()
	if err != nil {
		return
	}
}
