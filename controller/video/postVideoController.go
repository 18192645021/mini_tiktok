package video

import (
	"github.com/gin-gonic/gin"
	"mock_douyin_project/config"
	"mock_douyin_project/dao"
	"mock_douyin_project/service/video"
	"mock_douyin_project/util"
	"net/http"
	"path/filepath"
	"strconv"
)

var (
	suffixMap = map[string]struct{}{
		".mp4": {},
		".wmv": {},
		".mov": {},
		".avi": {},
	}
)

// PostVideo 上传视频接口
func PostVideo(c *gin.Context) {

	// 获取当前登录用户的userid
	idStr, _ := c.Get("user_id")
	userId := idStr.(int64)

	// 获取视频title
	title := c.PostForm("title")
	multipartForm, err := c.MultipartForm()
	if err != nil {
		PublishVideoError(c, err.Error())
		return
	}
	videoFile := multipartForm.File["data"]
	for _, file := range videoFile {
		// 校验文件格式是否正确
		suffix := filepath.Ext(file.Filename)

		if _, isExist := suffixMap[suffix]; !isExist {
			PublishVideoError(c, "视频格式不支持")
			continue
		}
		// 根据当前用户名生成文件名
		var videoCount int64
		if err := dao.NewVideoDAO().GetVideoCountByUserId(userId, &videoCount); err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		videoName := strconv.FormatInt(userId, 10) + "-" + strconv.FormatInt(videoCount, 10)
		savePath := filepath.Join(config.Global.StaticSourcePath, videoName+suffix)
		// 将视频保存到本地
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		// 截取一张封面图
		v2i := util.NewVideo2Image()
		util.ChangeVideoDefaultSuffix(suffix)
		v2i.InputPath = filepath.Join(config.Global.StaticSourcePath, videoName+util.GetDefaultVideoSuffix())
		v2i.OutputPath = filepath.Join(config.Global.StaticSourcePath, videoName+util.GetDefaultImageSuffix())
		v2i.FrameCount = 1
		res, err := v2i.GetQueryString()
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		if erro := v2i.ExecCommand(res); erro != nil {
			PublishVideoError(c, erro.Error())
			continue
		}
		// 将视频信息保存到数据库中
		if err := video.PublishVideo(userId, videoName+suffix, videoName+util.GetDefaultImageSuffix(), title); err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		PublishVideoSuccess(c, "视频上传成功")

	}

}
func PublishVideoSuccess(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, dao.CommonResponse{StatusCode: 0, StatusMsg: msg})
}
func PublishVideoError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, dao.CommonResponse{StatusCode: 1,
		StatusMsg: msg})
}
