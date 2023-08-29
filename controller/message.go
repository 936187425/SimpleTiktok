package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tiktok/model"
	"time"

	"github.com/gin-gonic/gin"
)

var tempChat = map[string][]model.Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	model.Response
	MessageList []model.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
// 发送消息
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
	username, err := model.RedisHandle.Get(token).Result()
	if err != nil {
		log.Printf(" MessageAction API: Token:%s donesn't exist\n", token)
		c.JSON(http.StatusForbidden, UserListResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Token:%s does not exist\n", token),
			},
			UserList: nil,
		})
		c.Abort()
		return
	}

	// 拿到当前登录的user
	fromUser := new(model.UserModel)
	model.MysqlHandle.Where("name = ?", username).First(fromUser)

	// 查找接收消息的用户
	toUserIdint, err := strconv.Atoi(toUserId)
	toUser := new(model.UserModel)
	result := model.MysqlHandle.Where("id = ?", toUserIdint).First(toUser)

	// 接收消息的用户不存在
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: -1,
			StatusMsg:  "fromUserId is not a correct id",
		})
		c.Abort()
		return
	}

	//插入表中
	msg := model.MessageModel{
		ToUserId:   toUserIdint,
		FromUserId: int(fromUser.Id),
		Content:    content,
		CreateTime: time.Now().Unix(),
	}
	model.MysqlHandle.Table("message_models").Create(&msg)
	c.JSON(http.StatusOK,
		model.Response{
			StatusCode: 0,
			StatusMsg:  "insert successfully!",
		},
	)
}

// MessageChat all users have same follow list
// 获得聊天记录
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	//先查询redis看看token是否存在
	username, err := model.RedisHandle.Get(token).Result()
	if err != nil {
		log.Printf(" MessageChat API: Token:%s donesn't exist\n", token)
		c.JSON(http.StatusForbidden, UserListResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Token:%s does not exist\n", token),
			},
			UserList: nil,
		})
		c.Abort()
		return
	}
	// token存在
	// 拿到当前登录的user
	fromUser := new(model.UserModel)
	model.MysqlHandle.Where("name = ?", username).First(fromUser)

	var msgContent []model.MessageModel
	model.MysqlHandle.Table("message_models").Where("from_user_id = ? and to_user_id = ?", fromUser.Id, toUserId).Find(&msgContent)

	var msgContentJson []model.Message
	for _, m := range msgContent {
		msg := model.Message{
			Id:         int64(m.Id),
			Content:    m.Content,
			CreateTime: strconv.FormatInt(m.CreateTime, 10),
		}
		msgContentJson = append(msgContentJson, msg)
	}

	c.JSON(http.StatusOK,
		ChatResponse{
			model.Response{
				StatusCode: 0,
			},
			msgContentJson,
		})
}
