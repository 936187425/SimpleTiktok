package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"strconv"
	"tiktok/model"
)

// FavoriteAction no practical effect, just check if token is valid
// 点赞后，添加到favorite表中
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	// 视频id
	video_id := c.Query("video_id")
	videoId, err := strconv.ParseInt(video_id, 10, 64)
	// 1点赞 2取消点赞
	action_type := c.Query("action_type")
	action, err := strconv.Atoi(action_type)

	// 验证token
	username, err := model.RedisHandle.Get(token).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "token已失效"})
		c.Abort()
		return
	}

	// 拿到当前登录的user
	userNow := new(model.UserModel)
	model.MysqlHandle.Where("name = ?", username).First(userNow)

	// 拿到当前video对象
	video := new(model.VideoModel)
	result := model.MysqlHandle.Where("id = ?", videoId).First(video)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "该视频不存在"})
		return
	}

	// 拿到当前video的user
	userAuthor := new(model.UserModel)
	model.MysqlHandle.Where("id = ?", video.UserID).First(userAuthor)

	favorite := new(model.Favorite)
	result = model.MysqlHandle.Where("user_id = ? and video_id = ?", userNow.Id, videoId).First(favorite)
	if action == 1 { // 点赞，先看一下该视频是否已经点过赞
		if result.RowsAffected != 0 {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "该视频已点赞，无法重复点赞"})
		} else {
			// favorite表插入记录
			favorite.UserId = userNow.Id
			favorite.VideoId = videoId
			model.MysqlHandle.Create(favorite)
			// 将该user的favorite_count字段+1
			userNow.FavoriteCount++
			model.MysqlHandle.Select("favorite_count").Save(userNow)
			// 将该video的favorite_count字段+1
			video.FavoriteCount++
			model.MysqlHandle.Select("favorite_count").Save(video)
			// 将该视频作者的total_favorited字段+1
			userAuthor.TotalFavorited++
			model.MysqlHandle.Select("total_favorited").Save(userAuthor)
			c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: "点赞成功"})
		}
	} else { // 取消点赞，先看一下数据库是否有这条点赞记录
		if result.RowsAffected == 0 {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "该视频还未点赞，无法取消点赞"})
		} else {
			// 在favorite表删除
			model.MysqlHandle.Delete(favorite)
			// 当前用户favorite_count字段-1
			userNow.FavoriteCount--
			model.MysqlHandle.Select("favorite_count").Save(userNow)
			// 将该video的favorite_count字段-1
			video.FavoriteCount--
			model.MysqlHandle.Select("favorite_count").Save(video)
			// 将该视频作者的total_favorited字段-1
			userAuthor.TotalFavorited--
			model.MysqlHandle.Select("total_favorited").Save(userAuthor)
			c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: "取消点赞成功"})
		}
	}
}

// FavoriteList all users have same favorite video list
// 根据user_id去favorite表中查询
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
