package dao

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
)

type Favorite struct {
	UserInfoId int64    `json:"user_info_id,omitempty" gorm:"many2many:user_favor_videos;"`
	VideoId    int64    `json:"video_id,omitempty" gorm:"many2many:user_favor_videos;"`
	User       UserInfo `json:"user,omitempty" gorm:"-"`
}

type FavoriteDAO struct {
}

var (
	favoriteDAO     *FavoriteDAO
	favoriteDaoOnce sync.Once
)

func NewFavoriteDao() *FavoriteDAO {
	favoriteDaoOnce.Do(func() {
		favoriteDAO = new(FavoriteDAO)
	})
	return favoriteDAO
}
func (*FavoriteDAO) IsUserAndVideoExistById(userId int64, videoId int64) bool {
	return videoDAO.IsVideoExistById(videoId) && userInfoDAO.IsUserExistById(userId)
}

func (*FavoriteDAO) FavoriteAction(userId int64, videoId int64) error {
	if favoriteDAO.IsUserAndVideoExistById(userId, videoId) {
		return DB.Transaction(func(tx *gorm.DB) error {
			var num int64
			DB.Where("userId=? AND videoId=?", userId, videoId).Count(&num)
			if num == 0 {
				err := tx.Table("user_favor_videos").Create(Favorite{UserInfoId: userId, VideoId: videoId}).Error
				if err != nil {
					log.Println(err)
				} else {
					err = tx.Exec("UPDATE videos v SET v.favorite_count = v.favorite_count+1 WHERE v.id=?", videoId).Error
					if err != nil {
						log.Println(err)
					}
				}
				return err
			} else {
				return errors.New("视频已被点赞，请刷新重试")
			}
		})
	}
	return errors.New("视频不存在")
}

func (*FavoriteDAO) FavoriteActionCancel(userId int64, videoId int64) error {
	if favoriteDAO.IsUserAndVideoExistById(userId, videoId) {
		return DB.Transaction(func(tx *gorm.DB) error {
			var num int64
			DB.Where("userId=? AND videoId=?", userId, videoId).Count(&num)
			if num != 0 {
				err := tx.Exec("DELETE FROM user_favor_videos ufv WHERE ufv.user_info_id=? AND ufv.video_id=? ", userId, videoId).Error
				if err != nil {
					log.Println(err)
					return err
				}
				err = tx.Exec("UPDATE videos v SET v.favorite_count = v.favorite_count-1 WHERE v.id=?", videoId).Error
				if err != nil {
					log.Println(err)
					return err
				}
				return nil
			} else {
				return errors.New("视频未被点赞，请刷新重试")
			}
		})
	}
	return errors.New("视频不存在")
}

func (*FavoriteDAO) QueryFavoriteList(userId int64) ([]*Video, error) {
	if NewUserInfoDAO().IsUserExistById(userId) {
		var VideoList []*Video
		err := DB.Raw("SELECT a.* from videos AS a inner join user_favor_videos AS b ON a.id = b.video_id AND b.user_info_id = ?", userId).Find(&VideoList).Error
		log.Println(err)
		return VideoList, err
	} else {
		err := errors.New("用户不存在")
		log.Println(err)
		return nil, err
	}
}
