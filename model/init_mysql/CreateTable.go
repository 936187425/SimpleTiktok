package init_mysql

import (
	"fmt"
	"os"
	"strings"
	"tiktok/model"
)

var userTable = "./Model.sql"

func CreateTable() {
	content, err := os.ReadFile(userTable)
	if err != nil {
		panic(fmt.Errorf("找不到UserModel配置文件！"))
	}
	// 将SQL内容拆分成多个SQL语句，逐句执行
	sqlLanguages := strings.Split(string(content), ";")
	for _, sql_ := range sqlLanguages {
		// 忽略空语句
		if strings.TrimSpace(sql_) == "" {
			continue
		}
		err := model.MysqlHandle.Exec(sql_)
		if err != nil {
			fmt.Printf("执行SQL语句时出现错误：%v\n", err)
			continue
		} else {
			fmt.Printf("SQL语句执行成功")
		}
	}
}
