package dao

import (
	"gorm.io/gorm"
	"log"
	"sync"
)

type UserRelation struct {
	UserInfoId int64 `gorm:"primaryKey;column:user_info_id"`
	FollowId   int64 `gorm:"primaryKey;column:follow_id"`
}

type UserRelationDao struct {
}

var (
	userRelationDao  *UserRelationDao
	userRelationOnce sync.Once
)

func NewUserRelationDao() *UserRelationDao {
	userRelationOnce.Do(func() {
		userRelationDao = new(UserRelationDao)
	})
	return userRelationDao
}

func (u *UserRelationDao) AddUserRelation(userRelation *UserRelation, userId int64, toUserId int64) error {
	if userRelation == nil || userId == 0 || toUserId == 0 {
		return ErrIvdPtr
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&userRelation).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE `user_infos` SET follow_count = follow_count + 1 WHERE id = ?", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE `user_infos` SET follower_count = follower_count + 1 WHERE id = ?", toUserId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (u *UserRelationDao) DeleteUserRelation(userRelation *UserRelation, userId int64, toUserId int64) error {
	if userRelation == nil || userId == 0 || toUserId == 0 {
		return ErrIvdPtr
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&userRelation).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE `user_infos` SET follow_count = follow_count - 1 WHERE id = ?", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE `user_infos` SET follower_count = follower_count - 1 WHERE id = ?", toUserId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (u *UserRelationDao) IsUserRelationExistByTwoId(userId int64, followId int64) bool {
	var userRelation UserRelation
	if err := DB.Where("user_info_id = ? AND follow_id = ?", userId, followId).First(&userRelation).Error; err != nil {
		log.Println(err)
	}
	if userRelation.UserInfoId == 0 {
		return false
	}
	return true
}

func (u *UserRelationDao) GetFollowRelationList(userId int64, userRelationList *[]UserRelation) error {
	return DB.Where("user_info_id = ?", userId).Find(userRelationList).Error
}

func (u *UserRelationDao) GetFollowerRelationList(followId int64, userRelationList *[]UserRelation) error {
	return DB.Where("follow_id = ?", followId).Find(userRelationList).Error
}

func (u *UserRelationDao) GetFriendRelationList(userId int64, followId, userRelationList *[]UserRelation) error {
	return DB.Where("user_info_id = ? AND follow_id = ?", userId, followId).Find(userRelationList).Error
}
