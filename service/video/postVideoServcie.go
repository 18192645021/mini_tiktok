package video

import (
	"fmt"
	"mock_douyin_project/config"
	"mock_douyin_project/dao"
)

// PublishVideo 投稿视频
func PublishVideo(userId int64, videoName, coverName, title string) error {
	return NewPostVideoFlow(userId, videoName, coverName, title).Do()
}

type PostVideoFlow struct {
	videoName string
	coverName string
	title     string
	userId    int64

	video *dao.Video
}

func NewPostVideoFlow(userId int64, videoName, coverName, title string) *PostVideoFlow {
	return &PostVideoFlow{
		videoName: videoName,
		coverName: coverName,
		userId:    userId,
		title:     title,
	}
}

func (p *PostVideoFlow) Do() error {
	p.videoName = fmt.Sprintf("http://%s:%d/static/%s", config.Global.IP, config.Global.Port, p.videoName)
	p.coverName = fmt.Sprintf("http://%s:%d/static/%s", config.Global.IP, config.Global.Port, p.coverName)
	video := &dao.Video{
		UserInfoId: p.userId,
		PlayUrl:    p.videoName,
		CoverUrl:   p.coverName,
		Title:      p.title,
	}
	if err := dao.NewVideoDAO().InsertVideo(video); err != nil {
		return err
	}
	return nil

}
