package dao

import (
	"sync"
)

type Message struct {
	Id         int64  `json:"id" gorm:"primary_key; auto_increment; not null"`
	ToUserId   int64  `json:"to_user_id" gorm:"to_user_id"`
	FromUserId int64  `json:"from_user_id" gorm:"from_user_id"`
	Content    string `json:"content" gorm:"content"`
	CreateTime int64  `json:"create_time" gorm:"create_time"`
}

type MessageDao struct {
}

var (
	messageDao     *MessageDao
	messageDaoOnce sync.Once
)

func NewMessageDao() *MessageDao {
	messageDaoOnce.Do(func() {
		messageDao = new(MessageDao)
	})
	return messageDao
}

func (m *MessageDao) AddMessage(newMessage *Message) error {
	if newMessage == nil {
		return ErrIvdPtr
	}
	return DB.Create(&newMessage).Error
}

func (m *MessageDao) SelectMessageList(userId int64, toUserId int64, preMsgTime int64, messageList *[]Message) error {
	return DB.
		Where(DB.Where("from_user_id = ? AND to_user_id = ?", userId, toUserId).
			Or("from_user_id = ? AND to_user_id = ?", toUserId, userId)).
		Where("create_time > ?", preMsgTime).
		Order("create_time asc").
		Find(&messageList).Error
}

func (m *MessageDao) SelectLastMessage(userId int64, toUserId int64, lastMessage *Message) error {
	//return DB.
	//	Where(DB.Where("from_user_id = ? AND to_user_id = ?", userId, toUserId).
	//		Or("from_user_id = ? AND to_user_id = ?", toUserId, userId)).
	//	Last(&lastMessage).Error

	row := DB.
		Where(DB.Where("from_user_id = ? AND to_user_id = ?", userId, toUserId).
			Or("from_user_id = ? AND to_user_id = ?", toUserId, userId))
	var num int64
	row.Count(&num)
	if num == 0 {
		return nil
	} else {
		lastMessage = nil
		return row.Last(&lastMessage).Error
	}
}
