package video

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mock_douyin_project/dao"
	midddleware "mock_douyin_project/middleware"
	"mock_douyin_project/service/video"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	dao.CommonResponse
	*video.FeedVideoList
}

func FeedVideoList(c *gin.Context) {
	p := NewProxyFeedVideoList(c)
	// 从token中拿到用户id
	token, ok := c.GetQuery("token")
	// 判断是否为登录状态
	if !ok {
		err := p.DoNoToken()
		if err != nil {
			p.FeedVideoListError(err.Error())
		}
		return
	}
	err := p.DoHasToken(token)
	if err != nil {
		p.FeedVideoListError(err.Error())
	}

}

type ProxyFeedVideoList struct {
	*gin.Context
}

func NewProxyFeedVideoList(c *gin.Context) *ProxyFeedVideoList {
	return &ProxyFeedVideoList{Context: c}
}

// DoNoToken 未登录状态获取视频流
func (p *ProxyFeedVideoList) DoNoToken() error {
	// 获取最近一次视频流中最早的时间
	latestTime := p.ParseLatestTime()
	videoList, err := video.QueryFeedVideoList(0, latestTime)
	if err != nil {
		return err
	}
	p.FeedVideoListSuccess(videoList)
	return nil
}

// DoHasToken 登录状态获取视频流
func (p *ProxyFeedVideoList) DoHasToken(token string) error {
	var claim *midddleware.Claims
	var err error
	if claim, err = midddleware.ParseToken(token); err == nil {
		if time.Now().Unix() > claim.ExpiresAt {
			// token超时
			return errors.New("token超时")
		}
		latestTime := p.ParseLatestTime()
		videoList, err := video.QueryFeedVideoList(claim.UserId, latestTime)
		if err != nil {
			return err
		}
		p.FeedVideoListSuccess(videoList)
		return nil
	}
	return err
}
func (p *ProxyFeedVideoList) ParseLatestTime() time.Time {
	rawTimestamp := p.Query("latest_time")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err == nil {
		latestTime = time.Unix(0, intTime*1e6)
	}
	return latestTime
}
func (p *ProxyFeedVideoList) FeedVideoListError(message string) {
	p.JSON(http.StatusOK, FeedResponse{CommonResponse: dao.CommonResponse{
		StatusCode: 1,
		StatusMsg:  message,
	}})
}
func (p *ProxyFeedVideoList) FeedVideoListSuccess(videoList *video.FeedVideoList) {
	p.JSON(http.StatusOK, FeedResponse{
		CommonResponse: dao.CommonResponse{
			StatusCode: 0,
		},
		FeedVideoList: videoList,
	})
}
