package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type server struct {
	Port int
	Host string
}

type mysql struct {
	Port     int
	Host     string
	User     string
	Password string
	Database string
}

type redis struct {
	Port int
	Host string
}

// Server 用于保存vipper读到的服务器配置文件
var Server *server

// Mysql 用于保存vipper读到的数据库配置文件
var Mysql *mysql

// Redis 用于保存vipper读到的Redis配置文件
var Redis *redis

func init() {
	Server = new(server)
	Mysql = new(mysql)
	Redis = new(redis)
	viper.SetConfigName("tiktok") // 设置配置文件名
	viper.SetConfigType("toml")   // 设置配置文件后缀
	viper.AddConfigPath(".")      // 设置配置文件路径
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("配置文件读取错误！"))
	}

	Server.Port = viper.GetInt("server.port")
	Server.Host = viper.GetString("server.host")

	Mysql.Port = viper.GetInt("mysql.port")
	Mysql.Host = viper.GetString("mysql.host")
	Mysql.User = viper.GetString("mysql.user")
	Mysql.Password = viper.GetString("mysql.password")
	Mysql.Database = viper.GetString("mysql.database")

	Redis.Port = viper.GetInt("redis.port")
	Redis.Host = viper.GetString("redis.host")
}

func (m *mysql) DSN() string {
	//return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=local", m.User, m.Password, m.Host, m.Port, m.Database)
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", m.User, m.Password, m.Host, m.Port, m.Database)
}
