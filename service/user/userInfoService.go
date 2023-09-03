package user

import (
	"mock_douyin_project/dao"
)

func DoQueryUserInfoByUserId(userId int64, info *dao.UserInfo) error {

	userinfoDao := dao.NewUserInfoDAO()
	var userInfo dao.UserInfo
	err := userinfoDao.QueryUserInfoById(userId, &userInfo)
	if err != nil {
		return err
	}
	info = &userInfo
	return nil
}
