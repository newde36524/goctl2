//nolint:dupl
//nolint:staticcheck
//nolint:errchkjson
//nolint:errcheck
package testlogic

import (
	"context"
	{{if .HasResp}}"encoding/json"{{end}}
	"os"
	"path/filepath"
	"testing"

	"app-server/api/{{.ProjectName}}/internal/config"
	{{if .NoTypes}}"app-server/api/{{.ProjectName}}/internal/types"{{end}}

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"

	{{.ImportPackages}}
)

// nolint: staticcheck
func Test_{{.function}}(t *testing.T) {
	getFilePath := func(rootDir string, subPath ...string) string {
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
	{{if .HasResp}}toJson := func(resp {{.ResponseType}},err error) (string, error){
		bs, _ := json.Marshal(resp) //nolint:errchkjson
		return string(bs), err
	}{{end}}
	var configFile = getFilePath("{{.ProjectName}}", "etc", "{{.ProjectName}}-dev.yaml")
	var c config.Config
	conf.MustLoad(configFile, &c)
	svcCtx := svc.NewServiceContext(c)

	ctx := context.WithValue(context.Background(), "userId", int64(1101001))
	l := {{.LogicName}}.New{{.logic}}(ctx, svcCtx)
	{{if .HasResp}}resp ,{{end}}err := l.{{.Call}}({{if .HasRequest}}&types.{{.RequestType}}{}{{end}})
	{{if .HasResp}}logx.Debug(toJson(resp ,err)){{else}}if err != nil {
		logx.Error(err)
	}{{end}}
}
