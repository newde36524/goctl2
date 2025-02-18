//nolint:dupl
//nolint:staticcheck
//nolint:errchkjson
//nolint:errcheck
package testlogic

import (
	"context"
	{{if .hasReply}}"encoding/json"{{end}}
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"app-server/service/sys/internal/config"
	{{.serviceName}}logic "app-server/service/{{.ProjectName}}/internal/logic/{{.serviceName}}"
	{{.imports}}

	"github.com/zeromicro/go-zero/core/conf"
)

// nolint: staticcheck
func Test_{{.method}}(t *testing.T) {
	var (
		getFilePath = func(rootDir string, subPath ...string) string {
			base := os.Args[0]
			for {
				base = filepath.Dir(base)
				if filepath.Base(base) == rootDir {
					elems := []string{base}
					return filepath.Join(append(elems, subPath...)...)
				}
				if len(filepath.Base(base)) == 0 {
					return ""
				}
			}
		}
		{{if .hasReply}}toJson = func(resp *{{.responseType}}) (string, error){
			bs, err := json.Marshal(resp) //nolint:errchkjson
			return string(bs), err
		}{{end}}
		configFile = getFilePath("{{.ProjectName}}", "etc", "{{.ProjectName}}-dev.yaml")
		c          config.Config
	)
	conf.MustLoad(configFile, &c)
	ctx, svcCtx := context.WithValue(context.Background(), "userId", int64(1101001)), svc.NewServiceContext(c)

	l := {{.serviceName}}logic.New{{.logicName}}(ctx, svcCtx)
	{{if .hasReply}}resp ,{{end}}err := l.{{.method}}({{if .hasReq}}&{{.request}}{}{{end}})
	{{if .hasReply}}fmt.Println(toJson(resp))
	if err != nil {
		fmt.Println(err)
		t.Failed()
	}{{else}}if err != nil {
		fmt.Println(err)
		t.Failed()
	}{{end}}
}
