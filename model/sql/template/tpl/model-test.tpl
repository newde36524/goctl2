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
)

var {{.lowerStartCamelObject}}Model = model.New{{.upperStartCamelObject}}Model(conn)

func Test_{{.upperStartCamelObject}}_Insert(t *testing.T) {
	_, err := {{.lowerStartCamelObject}}Model.Insert(context.TODO(), nil, &model.{{.upperStartCamelObject}}{{.feilds}})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_{{.upperStartCamelObject}}_Delete(t *testing.T) {
	err := {{.lowerStartCamelObject}}Model.Delete(context.TODO(), nil, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_{{.upperStartCamelObject}}_Update(t *testing.T) {
	_, err := {{.lowerStartCamelObject}}Model.Update(context.TODO(), nil, &model.{{.upperStartCamelObject}}{{.feilds}})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_{{.upperStartCamelObject}}_Select(t *testing.T) {
	v, err := {{.lowerStartCamelObject}}Model.FindOne(context.TODO(), 1)
	if err != nil {
		t.Fatal(err)
	}
	bs, _ := json.Marshal(v)
	fmt.Println(string(bs))
}
