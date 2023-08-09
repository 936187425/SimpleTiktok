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

// tokenTime token 的失效时间
var tokenTime = time.Minute * 100

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]model.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

//var userIdSequence = int64(1)

type UserLoginResponse struct {
	model.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	model.Response
	User model.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	u := new(model.UserModel)
	// 写入数据库
	u.Name = username
	u.Password = password
	res := model.MysqlHandle.Create(u)
	if res.Error != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "该用户已被注册！"},
			UserId:   0,
			Token:    "",
		})
		c.Abort()
		return
	}
	err := model.RedisHandle.Set(token, username, tokenTime).Err()
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: fmt.Sprintf("连接redis出错，错误信息：%v", err)},
			UserId:   0,
			Token:    "",
		})
		c.Abort()
		return
	}
	fmt.Printf(model.RedisHandle.DBSize().String())
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: model.Response{StatusCode: 0, StatusMsg: "注册成功！"},
		UserId:   u.Id,
		Token:    token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	u := new(model.UserModel)
	model.MysqlHandle.Where("name = ?", username).First(u)
	if u.Id == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "该用户未被注册！"},
			UserId:   u.Id,
			Token:    "",
		})
		c.Abort()
		return
	}

	// 检查密码是否一致
	if u.Password != password {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "密码不正确，请重新输入密码！"},
			UserId:   u.Id,
			Token:    "",
		})
		c.Abort()
		return
	} else {
		model.RedisHandle.Set(token, username, tokenTime)
		fmt.Printf(model.RedisHandle.DBSize().String())
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0, StatusMsg: "登录成功！"},
			UserId:   u.Id,
			Token:    token,
		})
	}

}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	// 验证token
	username, err := model.RedisHandle.Get(token).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "token 已失效！"},
			User:     model.User{}})
		c.Abort()
		return
	}
	// 根据token查找用户
	u := new(model.UserModel)
	user_id := c.Query("user_id")
	model.MysqlHandle.Where("name = ?", username).First(u)
	if strconv.Itoa(int(u.Id)) != user_id {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "身份验证不正确！"},
			User:     model.User{},
		})
		c.Abort()
		return
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 0, StatusMsg: "查询成功"},
			User: model.User{
				Id:              u.Id,
				Name:            u.Name,
				FollowCount:     0,
				FollowerCount:   0,
				IsFollow:        false,
				Avatar:          u.Avatar,
				BackgroundImage: u.BackgroundImage,
				Signature:       "gogogo",
				TotalFavorited:  0,
				WorkCount:       u.WorkCount,
				FavoriteCount:   0,
			},
		})
	}
}
