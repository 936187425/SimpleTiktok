package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"strconv"
	"tiktok/model"
	"time"
)

type CommentListResponse struct {
	model.Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	model.Response
	Comment model.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	videoId, _ := strconv.Atoi(video_id)
	actionType := c.Query("action_type") // 如果为1，表示评论，如果为2，表示删除评论

	// 验证token
	username, err := model.RedisHandle.Get(token).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "token已失效"})
		c.Abort()
		return
	}

	// 拿到当前登录的user
	var user model.User
	model.MysqlHandle.Table("user_models").Where("name = ?", username).First(&user)
	fmt.Println(user)

	// 拿到评论的视频
	video := new(model.VideoModel)
	result := model.MysqlHandle.Table("video_models").Where("id = ?", videoId).First(video)
	fmt.Println(video)

	// 没有该视频
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "Video doesn't exist"})
		return
	}

	if actionType == strconv.FormatInt(1, 10) { // 评论
		text := c.Query("comment_text")
		if text == "" {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "没有输入评论，无法评论"})
			return
		}

		// 评论表中添加一条记录
		comment := new(model.CommentModel)
		comment.UserId = user.Id
		comment.VideoId = int64(videoId)
		comment.Content = text
		currentTime := time.Now()
		comment.CreateDate = currentTime
		model.MysqlHandle.Create(comment)
		// 该视频的评论数+1
		video.CommentCount++
		model.MysqlHandle.Select("comment_count").Save(video)

		c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0, StatusMsg: "评论成功"},
			Comment: model.Comment{
				Id:         comment.Id,
				User:       user,
				Content:    text,
				CreateDate: currentTime.Format("2006-01-02 15:04:05"),
			}})
	} else if actionType == strconv.FormatInt(2, 10) { // 删除评论
		comment_id := c.Query("comment_id")
		if comment_id == "" {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "评论id为空，无法删除"})
			return
		}

		commentId, _ := strconv.Atoi(comment_id)
		// 查看是否有该评论
		comment := new(model.CommentModel)
		result := model.MysqlHandle.Where("id = ?", commentId).First(comment)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "该评论不存在，删除失败"})
			return
		}

		// 删除评论
		model.MysqlHandle.Delete(comment, commentId)
		// 该视频评论数-1
		video.CommentCount--
		model.MysqlHandle.Select("comment_count").Save(video)

		c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0, StatusMsg: "删除评论成功"},
			Comment: model.Comment{
				Id:         comment.Id,
				User:       user,
				Content:    comment.Content,
				CreateDate: comment.CreateDate.Format("2006-01-02 15:04:05"),
			}})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	videoId, _ := strconv.Atoi(video_id)

	// 验证token
	_, err := model.RedisHandle.Get(token).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "token已失效"})
		c.Abort()
		return
	}

	// 拿到该视频
	video := new(model.VideoModel)
	result := model.MysqlHandle.Table("video_models").Where("id = ?", videoId).First(video)

	// 视频不存在
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    model.Response{StatusCode: 1, StatusMsg: "该视频不存在"},
			CommentList: nil,
		})
		return
	}

	// 根据videoId查询其评论
	var commonList []model.CommentModel
	model.MysqlHandle.Where("video_id = ?", videoId).Order("create_date DESC").Find(&commonList)

	if len(commonList) == 0 {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    model.Response{StatusCode: 1, StatusMsg: "该视频目前没有评论"},
			CommentList: nil,
		})
		return
	}

	var commentJson []model.Comment
	for _, comment := range commonList {
		t_comment := model.Comment{}
		t_comment.Id = comment.Id
		user := model.User{}
		model.MysqlHandle.Table("user_models").Where("id = ?", comment.UserId).First(&user)
		t_comment.User = user
		t_comment.Content = comment.Content
		t_comment.CreateDate = comment.CreateDate.Format("2006-01-02 15:04:05")
		commentJson = append(commentJson, t_comment)
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    model.Response{StatusCode: 0, StatusMsg: "视频评论查询成功"},
		CommentList: commentJson,
	})
}
