package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mock_douyin_project/dao"
	midddleware "mock_douyin_project/middleware"
	"net/http"
)

type UserLoginResponse struct {
	dao.CommonResponse
	*LoginResponse
}

type LoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

const (
	MaxUsernameLength = 100
	MaxPasswordLength = 20
	MinPasswordLength = 8
)

func LoginService(c *gin.Context) {
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
	response, err := QueryUserLogin(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: dao.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: dao.CommonResponse{StatusCode: 0},
		LoginResponse:  response,
	})

}
func QueryUserLogin(username, password string) (*LoginResponse, error) {
	flow := QueryUserLoginFlow{username: username, password: password}
	return flow.Do()

}

type QueryUserLoginFlow struct {
	username string
	password string

	data   *LoginResponse
	userid int64
	token  string
}

func (q *QueryUserLoginFlow) Do() (*LoginResponse, error) {

	// 对参数进行验证
	if err := q.checkUsernameAndPassword(); err != nil {
		return nil, err
	}
	// 查询数据库准备数据
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	//打包最终数据
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.data, nil
}
func (q *QueryUserLoginFlow) checkUsernameAndPassword() error {
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
func (q *QueryUserLoginFlow) prepareData() error {
	userLoginDao := dao.NewUserLoginDao()
	var login dao.UserLogin
	err := userLoginDao.QueryUserLogin(q.username, q.password, &login)
	if err != nil {
		return err
	}
	q.userid = login.UserInfoId

	token, err := midddleware.GenerateToken(login)
	if err != nil {
		return err
	}
	q.token = token
	return nil
}

func (q *QueryUserLoginFlow) packData() error {
	q.data = &LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}
