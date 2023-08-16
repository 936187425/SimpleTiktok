package model

import "time"

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount uint   `json:"favorite_count,omitempty"`
	CommentCount  uint   `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title"`
}

type VideoModel struct {
	Id            int64     `json:"id,omitempty" db:"id" gorm:"not null;unique"`
	UserID        int64     `db:"user_id"`
	PlayUrl       string    `json:"play_url" json:"play_url,omitempty" db:"play_url"`
	CoverUrl      string    `json:"cover_url,omitempty" db:"cover_url"`
	FavoriteCount uint      `json:"favorite_count,omitempty" db:"favorite_count"`
	CommentCount  uint      `json:"comment_count,omitempty" db:"comment_count"`
	IsFavorite    bool      `json:"is_favorite,omitempty" db:"is_favorite"` // 当前登录用户对该视频是否点赞
	Title         string    `json:"title" db:"title" gorm:"not null"`
	Extension     string    `json:"extension" db:"extension"`
	CreatedTime   time.Time `json:"created_time" db:"created_time" gorm:"not null"`
	UpdatedTime   time.Time `json:"updated_time" db:"updated_time"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     uint   `json:"follow_count,omitempty"`
	FollowerCount   uint   `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  uint   `json:"total_favorited"`
	WorkCount       uint   `json:"work_count"`
	FavoriteCount   uint   `json:"favorite_count"`
}

type UserModel struct {
	Id              int64  `json:"id,omitempty" db:"id" gorm:"not null;unique"`
	Name            string `json:"name,omitempty" db:"username" gorm:"unique;not null"`
	Password        string `json:"password" db:"password" gorm:"not null"`
	FollowCount     uint   `json:"follow_count,omitempty" db:"follow_count"`
	FollowerCount   uint   `json:"follower_count,omitempty" db:"follower_count"`
	IsFollow        bool   `json:"is_follow,omitempty" db:"is_follow"`
	Avatar          string `json:"avatar" db:"avatar"`
	BackgroundImage string `json:"background_image" db:"background_image"`
	Signature       string `json:"signature" db:"signature"`
	TotalFavorited  uint   `json:"total_favorited" db:"total_favorited"`
	WorkCount       uint   `json:"work_count" db:"work_count"`
	FavoriteCount   uint   `json:"favorite_count" db:"favorite_count"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type Favorite struct {
	Id int64 `json:"id,omitempty" db:"id" gorm:"not null;unique"`
	//CreateTime time.Time `json:"createTime,omitempty" db:"create_time" gorm:"not null;default:current_timestamp"`
	//UpdateTime time.Time `json:"deleteTime,omitempty" db:"delete_time" gorm:"not null;default:current_timestamp"`
	UserId  int64 `json:"id,omitempty" db:"id" gorm:"not null;"`
	VideoId int64 `json:"id,omitempty" db:"id" gorm:"not null;"`
}
