package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mock_douyin_project/dao"
	midddleware "mock_douyin_project/middleware"
	"net/http"
)

type UserRegisterResponse struct {
	dao.CommonResponse
	*LoginResponse
}

func RegisterService(c *gin.Context) {
	// 获取请求中携带的用户名和密码
	username := c.Query("username")
	rawPwd, ok := c.Get("password")
	if !ok {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResponse: dao.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "密码解析出错",
			},
		})
		return
	}
	password := rawPwd.(string)
	response, err := NewUserLogin(username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResponse: dao.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		CommonResponse: dao.CommonResponse{StatusCode: 0},
		LoginResponse:  response,
	})
}

// NewUserLogin 注册用户使用的登录方法
func NewUserLogin(username, password string) (*LoginResponse, error) {
	flow := NewUserLoginFlow{username: username, password: password}
	return flow.Do()
}

type NewUserLoginFlow struct {
	username string
	password string
	data     *LoginResponse
	userid   int64
	token    string
}

func (q *NewUserLoginFlow) Do() (*LoginResponse, error) {
	// 对参数进行校验
	if err := q.checkUsernameAndPassword(); err != nil {
		return nil, err
	}
	// 将用户信息插入到数据库中
	if err := q.insertNewUser(); err != nil {
		return nil, err
	}
	// 返回响应信息
	if err := q.packResponse(); err != nil {
		return nil, err
	}
	return q.data, nil
}

// checkUsernameAndPassword 用户名和密码校验
func (q *NewUserLoginFlow) checkUsernameAndPassword() error {
	if q.username == "" {
		return errors.New("用户名为空")
	}
	if len(q.username) > MaxUsernameLength {
		return errors.New("用户名长度超出限制")
	}
	if q.password == "" {
		return errors.New("密码为空")
	}
	//if len(q.password) < MinPasswordLength || len(q.password) > MaxPasswordLength {
	//	return errors.New("密码长度不符合规定")
	//}
	return nil
}

// insertNewUser 插入新用户并生成token
func (q *NewUserLoginFlow) insertNewUser() error {

	// 准备用户信息
	userLogin := dao.UserLogin{Username: q.username, Password: q.password}
	userInfo := dao.UserInfo{User: &userLogin, Name: q.username}

	// 判断用户名是否已经存在
	userLoginDao := dao.NewUserLoginDao()
	if userLoginDao.IsUserExistByUsername(q.username) {
		return errors.New("用户名已存在")
	}
	// 执行插入操作
	userInfoDao := dao.NewUserInfoDAO()
	err := userInfoDao.AddUserInfo(&userInfo)
	if err != nil {
		return err
	}
	// 生成token
	token, err := midddleware.GenerateToken(userLogin)
	if err != nil {
		return nil
	}
	q.token = token
	q.userid = userInfo.Id
	return nil
}

func (q *NewUserLoginFlow) packResponse() error {
	q.data = &LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}
