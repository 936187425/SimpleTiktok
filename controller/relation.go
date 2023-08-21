package controller

import (
	"fmt"
	"log"
	"net/http"
	"tiktok/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	model.Response
	UserList []model.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: []model.User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: []model.User{DemoUser},
	})
}

// FriendList all users have same friend list
// 获取好友记录
func FriendList(c *gin.Context) {
	
	user_id:=c.Query("user_id")
	token:=c.Query("token")
	log.Printf("FriendLis API: User_id:%s Token:%s\n",user_id,token)

	//先查询redis看看token是否存在
	_,err:=model.RedisHandle.Get(token).Result()
	if err!=nil{
		log.Printf("FriendList API: Token:%s donesn't exist\n",token)
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
	//token是存在,查询数据库获得好友列表.
	var FriendListModel model.FriendListModel
	model.MysqlHandle.Table("friend_list_models").Where("name=?",user_id).First(&FriendListModel)
	log.Printf("FriendList API: Token:%s FriendList:%+v\n",token,FriendListModel)
	listFriendByte:=[]byte(FriendListModel.List)
	//反序列化
	var listSlice []model.User
	json.Unmarshal(listFriendByte,&listSlice)
	c.JSON(http.StatusOK,UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: listSlice, //listSlice为null说明该username没有好友.
	})
	c.Abort()
	return
}
