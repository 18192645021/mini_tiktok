package video

import (
	"errors"
	"mock_douyin_project/dao"
	"time"
)

const (
	MaxVideoNum = 30
)

type FeedVideoList struct {
	Videos   []*dao.Video `json:"video_list,omitempty"`
	NextTime int64        `json:"next_time,omitempty"`
}

func QueryFeedVideoList(userId int64, latestTime time.Time) (*FeedVideoList, error) {
	return NewQueryFeedVideoListFlow(userId, latestTime).Do()
}

type QueryFeedVideoListFlow struct {
	userId     int64
	latestTime time.Time

	videos    []*dao.Video
	nextTime  int64
	feedVideo *FeedVideoList
}

func NewQueryFeedVideoListFlow(userId int64, latestTime time.Time) *QueryFeedVideoListFlow {
	return &QueryFeedVideoListFlow{userId: userId, latestTime: latestTime}
}

func (q *QueryFeedVideoListFlow) Do() (*FeedVideoList, error) {
	// 如果latestTime参数为空，则将其设置为当前时间戳
	if q.latestTime.IsZero() {
		q.latestTime = time.Now()
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.feedVideo, nil
}

func (q *QueryFeedVideoListFlow) prepareData() error {
	err := dao.NewVideoDAO().GetVideoListByLimitAndTime(MaxVideoNum, q.latestTime, &q.videos)
	if err != nil {
		return err
	}
	// 用户如果为登录状态，则需要该用户是否点赞该视频的信息插入到每条videoInfo中
	// 得到当前返回的videoList中最早的视频的时间戳
	if err := q.FillVideoListFields(); err != nil {
		q.nextTime = time.Now().Unix()
	}
	return nil
}

func (q *QueryFeedVideoListFlow) packData() error {
	q.feedVideo = &FeedVideoList{
		Videos:   q.videos,
		NextTime: q.nextTime,
	}
	return nil
}
func (q *QueryFeedVideoListFlow) FillVideoListFields() error {
	if q.videos == nil || len(q.videos) == 0 {
		return errors.New("当前videos列表为空")
	}
	q.latestTime = q.videos[len(q.videos)-1].CreatedAt
	userInfoDao := dao.NewUserInfoDAO()
	// 添加video的作者信息
	for i := 0; i < len(q.videos); i++ {
		var userInfo dao.UserInfo
		if err := userInfoDao.QueryUserInfoById(q.videos[i].UserInfoId, &userInfo); err != nil {
			continue
		}
		if q.userId > 0 {
			// 填充当前登录用户和该视频用于的关注信息
			//userInfo.IsFollow = GetUserRelation(q.userId, userInfo.Id)
			// 填充当前登录用户是否点赞了该视频
			q.videos[i].IsFavorite = dao.NewVideoDAO().GetVideoIsFavoriteByUserIdAndVideoId(q.userId, q.videos[i].Id)
		}
		q.videos[i].Author = userInfo

	}
	return nil
}
