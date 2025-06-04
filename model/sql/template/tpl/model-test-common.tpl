package modeltest

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var conn = sqlx.NewMysql("root:123456@tcp(127.0.0.1:3306)/dev?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
