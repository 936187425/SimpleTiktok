package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"path/filepath"
	"strconv"
	"tiktok/config"
	"tiktok/model"
	"time"
)

var videoIdSequence = int64(1)

type VideoListResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.Query("token")
	title := c.Query("title")

	// 验证token
	username, err := model.RedisHandle.Get(token).Result()
	if err == redis.Nil {
		//fmt.Printf(err.Error())
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "token已失效"})
		c.Abort()
		return
	}

	// title不能为空
	if title == "" {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "标题不能为空！"})
		c.Abort()
		return
	}
	// 获取文件
	file, err_ := c.FormFile("data")
	if err_ != nil {
		fmt.Printf(err_.Error())
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "视频上传失败！"})
		c.Abort()
		return
	}
	// 获取文件扩展名
	ext := filepath.Ext(file.Filename)
	// 将视频信息记录到数据库
	u := new(model.UserModel)
	model.MysqlHandle.Where("name = ?", username).First(u)
	fmt.Printf(strconv.FormatInt(u.Id, 10))
	if u.Id == 0 {
		fmt.Printf(u.Name)
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "token 不正确！"})
		c.Abort()
		return
	}
	playurl := fmt.Sprintf("%s:%d/public/video/%d%s", config.Server.Host, config.Server.Port, videoIdSequence, ext)
	video := model.VideoModel{
		Id:          videoIdSequence,
		Title:       title,
		Extension:   ext,
		UserID:      u.Id,
		CreatedTime: time.Now(),
		PlayUrl:     playurl,
	}
	model.MysqlHandle.Create(&video)
	// 保存视频: public/video/id.ext
	path := fmt.Sprintf("public/video/%d%s", video.Id, ext)
	errvideo := c.SaveUploadedFile(file, "./"+path)
	videoIdSequence += 1
	if errvideo != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "视频保存失败，请重新上传！"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: "视频发布成功！"})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	// 验证token
	username, err := model.RedisHandle.Get(token).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "token已失效！"})
		c.Abort()
		return
	}
	// 查询所有video
	u := new(model.UserModel)
	model.MysqlHandle.Where("name = ?", username).First(u)
	if strconv.FormatInt(u.Id, 10) != userId {
		c.JSON(http.StatusOK, VideoListResponse{
			Response:  model.Response{StatusCode: 1, StatusMsg: "身份验证错误！"},
			VideoList: nil,
		})
		c.Abort()
		return
	}
	var videos []model.VideoModel
	model.MysqlHandle.Where("user_id = ?", u.Id).Find(&videos)
	//model.MysqlHandle.Where(&model.VideoModel{UserID: u.Id}).Find(videos)
	if len(videos) == 0 {
		c.JSON(http.StatusOK, VideoListResponse{
			Response:  model.Response{StatusCode: 1, StatusMsg: "该用户还未发布任何视频！"},
			VideoList: nil,
		})
		c.Abort()
		return
	}
	videoJson := make([]model.Video, 0, 30)
	for _, video := range videos {
		u := new(model.UserModel)
		model.MysqlHandle.Where("id = ?", video.UserID).First(u)
		user := model.User{
			Id:              u.Id,
			Name:            u.Name,
			FollowCount:     u.FollowCount,
			FollowerCount:   u.FollowerCount,
			IsFollow:        u.IsFollow,
			Avatar:          u.Avatar,
			BackgroundImage: u.BackgroundImage,
			Signature:       u.Signature,
			TotalFavorited:  u.TotalFavorited,
			WorkCount:       u.WorkCount,
			FavoriteCount:   u.FavoriteCount,
		}
		v := model.Video{
			Id:            video.Id,
			Author:        user,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		}
		videoJson = append(videoJson, v)
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  model.Response{StatusCode: 0, StatusMsg: "查询成功！"},
		VideoList: videoJson,
	})
}
