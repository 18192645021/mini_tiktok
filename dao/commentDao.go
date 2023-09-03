package dao

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Comment struct {
	Id         int64     `json:"id,omitempty"`
	UserInfoId int64     `json:"-"` //用于一对多关系的id
	VideoId    int64     `json:"-"` //一对多，视频对评论
	User       UserInfo  `json:"user,omitempty" gorm:"-"`
	Content    string    `json:"content,omitempty"`
	CreatedAt  time.Time `json:"-"`
	CreateDate string    `json:"create_date,omitempty" gorm:"-"`
}

type CommentDAO struct {
}

var (
	commentDAO     *CommentDAO
	commentDaoOnce sync.Once
)

func NewCommentDAO() *CommentDAO {
	commentDaoOnce.Do(func() {
		commentDAO = new(CommentDAO)
	})
	return commentDAO
}

func (c *CommentDAO) AddCommentAndUpdateCount(comment *Comment) error {
	if comment == nil {
		return errors.New("AddCommentAndUpdateCount comment空指针")
	}
	//执行事务
	return DB.Transaction(func(tx *gorm.DB) error {
		//添加评论数据
		if err := tx.Create(comment).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		//增加count
		if err := tx.Exec("UPDATE videos v SET v.comment_count = v.comment_count+1 WHERE v.id=?", comment.VideoId).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

func (c *CommentDAO) DeleteCommentAndUpdateCountById(commentId, videoId int64) error {
	//执行事务
	return DB.Transaction(func(tx *gorm.DB) error {
		//删除评论
		if err := tx.Exec("DELETE FROM comments WHERE id = ?", commentId).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		//减少count
		if err := tx.Exec("UPDATE videos v SET v.comment_count = v.comment_count-1 WHERE v.id=? AND v.comment_count>0", videoId).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}
func (*CommentDAO) IsCommentExistById(userId int64, videoId int64, commentId int64) bool {
	if !NewUserInfoDAO().IsUserExistById(userId) || !NewVideoDAO().IsVideoExistById(videoId) {
		return false
	} else {
		var comment Comment
		if err := DB.Where("id=?", commentId).First(&comment).Error; err != nil {
			log.Println(err)
			return false
		}
		if comment.Id == 0 {
			return false
		}
		return true
	}
}
func (*CommentDAO) QueryCommentListByVideoId(videoId int64) ([]*Comment, error) {
	if !NewVideoDAO().IsVideoExistById(videoId) {
		return nil, errors.New("视频不存在")
	} else {
		var commentList []*Comment
		err := DB.Raw("SELECT * from comments where video_id = ?  ", videoId).Find(&commentList).Error
		if err != nil {
			log.Println(err)
		}
		return commentList, err
	}
}
func (*CommentDAO) QueryCommentCountByVideoId(videoId int64, commentCount *int64) error {
	if !NewVideoDAO().IsVideoExistById(videoId) {
		return errors.New("视频不存在")
	}
	err := DB.Model(&Comment{}).Where("video_id=?", videoId).Count(commentCount).Error
	if err != nil {
		return errors.New("videoCommentCount error")
	}
	return nil
}
