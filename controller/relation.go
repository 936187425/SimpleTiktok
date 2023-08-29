package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"strconv"
	"tiktok/model"
	"time"
)

type UserListResponse struct {
	model.Response
	UserList []model.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")

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

	// 查找准备关注或取关的用户
	toUser := new(model.UserModel)
	result := model.MysqlHandle.Where("id = ?", to_user_id).First(toUser)
	// 关注或取关的用户不存在
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "关注或取关的用户不存在"})
		return
	}

	// 不能关注自己
	if to_user_id == strconv.FormatInt(userNow.Id, 10) {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "不能关注自己"})
		return
	}

	// 判断是否已经关注过该用户
	relation := new(model.Relation)
	result = model.MysqlHandle.Where("from_user_id = ? and to_user_id = ?", userNow.Id, to_user_id).First(relation)

	if action_type == "1" { // 关注操作
		// 已经关注过该用户
		if result.RowsAffected != 0 {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "该用户已关注，无需重复关注"})
			return
		}
		// 在relation表添加记录
		relation.FromUserId = userNow.Id
		relation.ToUserId = toUser.Id
		relation.CreateDate = time.Now()
		model.MysqlHandle.Create(relation)
		// 关注者的关注数+1
		userNow.FollowCount++
		model.MysqlHandle.Select("follow_count").Save(userNow)
		// 被关注者的粉丝数+1
		toUser.FollowerCount++
		model.MysqlHandle.Select("follower_count").Save(toUser)
		c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: "关注成功"})
	} else if action_type == "2" { // 取关操作
		// 没有关注过该用户
		if result.RowsAffected == 0 {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "未关注该用户，无法取关"})
			return
		}
		// 在relation表删除记录
		model.MysqlHandle.Delete(relation)
		// 关注者的关注数-1
		userNow.FollowCount--
		model.MysqlHandle.Select("follow_count").Save(userNow)
		// 被关注者的粉丝数-1
		toUser.FollowerCount--
		model.MysqlHandle.Select("follower_count").Save(toUser)
		c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: "取关成功"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")

	// 验证token
	_, err := model.RedisHandle.Get(token).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "token已失效"})
		c.Abort()
		return
	}

	// 查看该user是否存在
	user := new(model.UserModel)
	result := model.MysqlHandle.Where("id = ?", user_id).First(user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "查询的用户不存在"},
			UserList: nil,
		})
		return
	}

	// 在relation表中，根据user_id查询到所有该用户关注的记录
	var relations []model.Relation
	model.MysqlHandle.Where("from_user_id = ?", user_id).Find(&relations)

	if len(relations) == 0 {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "关注列表为空"},
			UserList: nil,
		})
		return
	}

	var follows []model.UserModel
	for _, relation := range relations {
		user := model.UserModel{}
		model.MysqlHandle.Table("user_models").Where("id = ?", relation.ToUserId).First(&user)
		follows = append(follows, user)
	}

	var followsJson []model.User
	for _, u := range follows {
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
		followsJson = append(followsJson, user)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "关注列表查询成功",
		},
		UserList: followsJson,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")

	// 验证token
	_, err := model.RedisHandle.Get(token).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "token已失效"})
		c.Abort()
		return
	}

	// 查看该user是否存在
	user := new(model.UserModel)
	result := model.MysqlHandle.Where("id = ?", user_id).First(user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "查询的用户不存在"},
			UserList: nil,
		})
		return
	}

	// 在relation表中，根据user_id查询到所有关注该用户的记录
	var relations []model.Relation
	model.MysqlHandle.Where("to_user_id = ?", user_id).Find(&relations)

	if len(relations) == 0 {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "粉丝列表为空"},
			UserList: nil,
		})
		return
	}

	var followers []model.UserModel
	for _, relation := range relations {
		user := model.UserModel{}
		model.MysqlHandle.Table("user_models").Where("id = ?", relation.FromUserId).First(&user)
		followers = append(followers, user)
	}

	var followersJson []model.User
	for _, u := range followers {
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
		followersJson = append(followersJson, user)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "粉丝列表查询成功",
		},
		UserList: followersJson,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	log.Printf("FriendLis API: User_id:%s Token:%s\n", user_id, token)

	//先查询redis看看token是否存在
	_, err := model.RedisHandle.Get(token).Result()
	if err != nil {
		log.Printf("FriendList API: Token:%s donesn't exist\n", token)
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

	// 判断user_id是否合法
	var user = model.UserModel{}
	result := model.MysqlHandle.Where("id = ?", user_id).First(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg:  "user_id不合法",
			},
			UserList: nil,
		})
		return
	}

	// 查找当前user的好友
	var friendList []model.UserModel

	// 查找当前user的关注列表
	var userFollowRelationList []model.Relation
	model.MysqlHandle.Where("from_user_id = ?", user.Id).Find(&userFollowRelationList)
	fmt.Println(userFollowRelationList)

	for _, relation := range userFollowRelationList {
		var t_relation = model.Relation{}
		// 如果对方也关注我了，那我们就是好友
		result = model.MysqlHandle.Where("from_user_id = ? and to_user_id = ?", relation.ToUserId, user_id).First(&t_relation)
		if result.RowsAffected == 0 {
			continue
		}
		var user = model.UserModel{}
		model.MysqlHandle.Where("id = ?", relation.ToUserId).First(&user)
		friendList = append(friendList, user)
	}

	var friendListJson []model.User
	for _, u := range friendList {
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
		friendListJson = append(friendListJson, user)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: friendListJson, //listSlice为null说明该username没有好友.
	})
}
