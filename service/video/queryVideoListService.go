package video

import (
	"errors"
	"mock_douyin_project/dao"
)

type List struct {
	Videos []*dao.Video `json:"video_list,omitempty"`
}

func QueryVideoListByUserId(userId int64) (*List, error) {
	var videoList List
	// 检查用户是否存在
	if !dao.NewUserInfoDAO().IsUserExistById(userId) {
		return &videoList, errors.New("用户不存在")
	}
	// 根据用户id查询投稿的视频
	var videos []*dao.Video
	err := dao.NewVideoDAO().GetVideoListByUserId(userId, &videos)
	if err != nil {
		return &videoList, err
	}
	// 查询作者信息
	var userinfo dao.UserInfo
	err = dao.NewUserInfoDAO().QueryUserInfoById(userId, &userinfo)
	if err != nil {
		return &videoList, err
	}
	// 填充信息(作者信息和自己的点赞情况)
	for i := range videos {
		videos[i].Author = userinfo
		// 用户自己是否对自己的视频点赞
		videos[i].IsFavorite = dao.NewVideoDAO().GetVideoIsFavoriteByUserIdAndVideoId(userId, videos[i].Id)
	}
	videoList.Videos = videos
	return &videoList, nil
}
