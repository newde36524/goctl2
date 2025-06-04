package modeltest

import (
	"context"
	{{if .sql}}"database/sql"{{end}}
	"encoding/json"
	"fmt"
	"{{.modelImports}}"
	"testing"
	{{if .time}}"time"
	{{end}}
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var {{.lowerStartCamelObject}}Conn = sqlx.NewMysql("root:123456@tcp(127.0.0.1:3306)/dev?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
var {{.lowerStartCamelObject}}Model = model.New{{.upperStartCamelObject}}Model({{.lowerStartCamelObject}}Conn)

func Test{{.upperStartCamelObject}}Insert(t *testing.T) {
	_, err := {{.lowerStartCamelObject}}Model.Insert(context.TODO(), nil, &model.{{.upperStartCamelObject}}{{.feilds}})
	if err != nil {
		t.Fatal(err)
	}
}

func Test{{.upperStartCamelObject}}Delete(t *testing.T) {
	err := {{.lowerStartCamelObject}}Model.Delete(context.TODO(), nil, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test{{.upperStartCamelObject}}Update(t *testing.T) {
	_, err := {{.lowerStartCamelObject}}Model.Update(context.TODO(), nil, &model.{{.upperStartCamelObject}}{{.feilds}})
	if err != nil {
		t.Fatal(err)
	}
}

func Test{{.upperStartCamelObject}}Select(t *testing.T) {
	v, err := {{.lowerStartCamelObject}}Model.FindOne(context.TODO(), 1)
	if err != nil {
		t.Fatal(err)
	}
	bs, _ := json.Marshal(v)
	fmt.Println(string(bs))
}
