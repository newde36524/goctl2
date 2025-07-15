package gogen

import (
	"fmt"
	"path"
	"strings"

	"github.com/newde36524/goctl2/api/spec"
	"github.com/newde36524/goctl2/config"
	"github.com/newde36524/goctl2/util/format"
)

func genApiTest(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	for _, g := range api.Service.Groups {
		for _, r := range g.Routes {
			err := genLogicByApiTest(dir, cfg, g, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genLogicByApiTest(dir string, cfg *config.Config, group spec.Group, route spec.Route) error {
	//todo 生成各个接口 路由 的api文件，包含 请求参数
	var (
		doc        = strings.TrimSpace(strings.TrimPrefix(getDoc(route.JoinedDoc()), "//"))
		httpMethod = route.Method
		routePath  = route.Path
	)

	var reqFeilds []string
	structType, ok := route.RequestType.(spec.DefineStruct)
	if ok {
		for _, member := range structType.Members {
			if len(member.Name) == 0 {
				continue
			}
			str := fmt.Sprintf("	\"%s\": %v,", lowerCamelCase(member.Name), GetTypeDefaultValue(member))
			reqFeilds = append(reqFeilds, str)
		}
	}
	var reqFeildsStr string
	if len(reqFeilds) > 0 {
		reqFeilds[len(reqFeilds)-1] = strings.TrimSuffix(reqFeilds[len(reqFeilds)-1], ",")
		reqFeildsStr = strings.Join(reqFeilds, "\n")
		reqFeildsStr = "\n" + reqFeildsStr + "\n"
	}
	reqFeildsStr = "{" + reqFeildsStr + "}"

	method, contentType := getHttpMethodAndContentType(httpMethod)
	text := fmt.Sprintf(`@HOST = http://127.0.0.1:1234
@authToken = {{login.response.body.data}}

# @name login
POST {{HOST}}/login
Content-Type: application/json

{
  "username":"admin",
  "password":"123456"
}
###

# %s
%s {{HOST}}%s
Content-Type: %s
Authorization: {{authToken}}

%s
###`, doc, method, routePath, contentType, reqFeildsStr)

	logic := getLogicName(route)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, logic)
	if err != nil {
		return err
	}
	filename := goFile + ".http"
	subDir := getLogicFolderPathApiTest(group, route)
	return genApiTestFile(dir, subDir, filename, text)
}

func getHttpMethodAndContentType(method string) (string, string) {
	switch method {
	case "get":
		return "GET", "application/x-www-form-urlencoded"
	case "post":
		return "POST", "application/json"
	}
	return strings.ToUpper(method), "application/json"
}

func getLogicFolderPathApiTest(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return apiDirTest
		}
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(apiDirTest, folder)
}
