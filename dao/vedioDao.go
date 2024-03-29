package dao

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Video struct {
	Id            int64       `json:"id,omitempty"`
	UserInfoId    int64       `json:"-"`
	Author        UserInfo    `json:"author,omitempty" gorm:"-"` //这里应该是作者对视频的一对多的关系，而不是视频对作者，故gorm不能存他，但json需要返回它
	PlayUrl       string      `json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url,omitempty"`
	FavoriteCount int64       `json:"favorite_count,omitempty"`
	CommentCount  int64       `json:"comment_count,omitempty"`
	IsFavorite    bool        `json:"is_favorite,omitempty"`
	Title         string      `json:"title,omitempty"`
	Users         []*UserInfo `json:"-" gorm:"many2many:user_favor_videos;"`
	Comments      []*Comment  `json:"-"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
}
type UserFavorVideo struct {
	UserInfoId int64 `gorm:"primaryKey"`
	VideoId    int64 `gorm:"primaryKey"`
}
type VideoDAO struct {
}

var (
	videoDAO  *VideoDAO
	videoOnce sync.Once
)

func NewVideoDAO() *VideoDAO {
	videoOnce.Do(func() {
		videoDAO = new(VideoDAO)
	})
	return videoDAO
}

func (v *VideoDAO) InsertVideo(video *Video) error {
	if video == nil {
		return errors.New("InsertVideo video nil")
	}
	return DB.Create(video).Error
}
func (v *VideoDAO) GetVideoListByUserId(userId int64, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByUserId videoList 空指针")
	}
	return DB.Where("user_info_id=?", userId).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title"}).
		Find(videoList).Error
}

func (v *VideoDAO) GetVideoIsFavoriteByUserIdAndVideoId(userid int64, videoid int64) bool {
	var userFavorVideo UserFavorVideo
	if err := DB.Where("user_info_id = ? AND video_id = ?", userid, videoid).First(&userFavorVideo); err != nil {
		log.Println(err)
	}
	if userFavorVideo.UserInfoId == 0 {
		return false
	}
	return true
}

func (v *VideoDAO) GetVideoListByLimitAndTime(limit int, latestTime time.Time, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("GetVideoListByLimitAndTime videoList 空指针")
	}
	return DB.Model(&Video{}).Where("created_at < ?", latestTime).
		Order("created_at ASC").Limit(limit).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title", "created_at", "updated_at"}).
		Find(videoList).Error
}
func (v *VideoDAO) GetVideoCountByUserId(userId int64, count *int64) error {
	if count == nil {
		return errors.New("QueryVideoCountByUserId count 空指针")
	}
	return DB.Model(&Video{}).Where("user_info_id=?", userId).Count(count).Error
}
func (v *VideoDAO) IsVideoExistById(videoId int64) bool {
	var video Video
	if err := DB.Where("Id = ?", videoId).First(&video); err != nil {
		log.Println(err)
	}
	if video.Id == 0 {
		return false
	}
	return true
}
