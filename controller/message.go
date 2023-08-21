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
	MessageList []model.MessageModel `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
// 发送消息
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
	result,err:=model.RedisHandle.Get(token).Result()
	if err!=nil{
		log.Printf(" MessageAction API: Token:%s donesn't exist\n",token)
		c.JSON(http.StatusForbidden,UserListResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg: fmt.Sprintf("Token:%s does not exist\n",token),
			},
			UserList:nil,
		})
		c.Abort()
		return
	}
	//token存在
	//先查询最大的id值
	var maxId int
	model.MysqlHandle.Table("message_models").Select("MAX(id)").Scan(&maxId)
	maxId++

	toUserIdint,err:=strconv.Atoi(toUserId)
	fromUserIdint,err:=strconv.Atoi(result)
	if err!=nil{
		c.JSON(http.StatusBadRequest,model.Response{
			StatusCode: -1,
			StatusMsg: "fromUserId is not a correct id",
		})
		c.Abort()
		return 
	}
	//插入表中
	msg:=model.MessageModel{
		Id:maxId,
		To_user_id: toUserIdint,
		From_user_id: fromUserIdint,
		Content: content,
		Create_time: time.Now().Unix(),
	}
	model.MysqlHandle.Table("message_models").Create(&msg)
	c.JSON(http.StatusOK,
		model.Response{
			StatusCode: 0,
			StatusMsg: "insert successfully!",
		},
	)
}

// MessageChat all users have same follow list
// 获得聊天记录
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	//先查询redis看看token是否存在
	result,err:=model.RedisHandle.Get(token).Result()
	if err!=nil{
		log.Printf(" MessageChat API: Token:%s donesn't exist\n",token)
		c.JSON(http.StatusForbidden,UserListResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg: fmt.Sprintf("Token:%s does not exist\n",token),
			},
			UserList:nil,
		})
		c.Abort()
		return
	}
	//token存在
	fromUserId:=result
	var msgContent []model.MessageModel
	model.MysqlHandle.Table("message_models").Where("from_user_id=? and to_user_id=?",fromUserId,toUserId).Find(&msgContent)

	c.JSON(http.StatusOK,
	ChatResponse{
		model.Response{
			StatusCode:0,
		},
		msgContent,
	})
}


