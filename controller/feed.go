package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/model"
	"time"
)

type FeedResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  time.Time     `json:"next_time,omitempty"`
}

var limitNum = 30

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	ltTime := c.Query("latest_time")
	const base_format = "2006-01-02 15:04:05"
	latestTime, _ := time.Parse(base_format, ltTime)

	now := time.Now()
	// 查询数据库中符合条件的视频
	var videos []model.VideoModel
	model.MysqlHandle.Where("created_time > ? and created_time < ?", latestTime, now).Order("created_time desc").Limit(limitNum).Find(&videos)

	if len(videos) == 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  model.Response{StatusCode: 1, StatusMsg: "暂无视频"},
			VideoList: nil,
			NextTime:  now,
		})
		c.Abort()
		return
	}
	videosJson := make([]model.Video, 0, limitNum)
	for _, video := range videos {
		u := new(model.User)
		model.MysqlHandle.Where("id = ?", video.UserID).First(u)
		v := model.Video{
			Id:            video.Id,
			Author:        *u,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         video.Title,
		}
		videosJson = append(videosJson, v)
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  model.Response{StatusCode: 0, StatusMsg: "视频流获取成功！"},
		VideoList: videosJson,
		NextTime:  videos[len(videos)-1].CreatedTime,
	})
}
