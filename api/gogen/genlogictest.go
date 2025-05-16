package gogen

import (
	_ "embed"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/newde36524/goctl2/api/spec"
	"github.com/newde36524/goctl2/config"
	"github.com/newde36524/goctl2/util"
	"github.com/newde36524/goctl2/util/format"
	"github.com/newde36524/goctl2/util/pathx"
)

//go:embed logic.tpl.test
var logicTestTemplate string

func genLogicTest(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	for _, g := range api.Service.Groups {
		for _, r := range g.Routes {
			err := genLogicByRouteTest(dir, rootPkg, cfg, g, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genLogicByRouteTest(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route) error {
	logic := getLogicName(route)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, logic)
	if err != nil {
		return err
	}

	imports := genLogicImports(route, rootPkg)
	var responseString string
	var returnString string
	var requestString string
	if len(route.ResponseTypeName()) > 0 {
		resp := responseGoTypeName(route, typesPacket)
		responseString = "(resp " + resp + ", err error)"
		returnString = "return"
	} else {
		responseString = "error"
		returnString = "return nil"
	}
	if len(route.RequestTypeName()) > 0 {
		requestString = "req *" + requestGoTypeName(route, typesPacket)
	}

	handler := getHandlerName(route)
	handlerPath := getHandlerFolderPath(group, route)
	pkgName := handlerPath[strings.LastIndex(handlerPath, "/")+1:]
	logicName := defaultLogicPackage
	if handlerPath != handlerDir {
		handler = strings.Title(handler)
		logicName = pkgName
	}
	importPackages := genHandlerImports(group, route, rootPkg)
	subDir := getLogicFolderPathTest(group, route)
	var reqFeilds []string
	structType, ok := route.RequestType.(spec.DefineStruct)
	if ok {
		for _, member := range structType.Members {
			if len(member.Name) == 0 {
				continue
			}
			str := fmt.Sprintf("%s: %v,", upperCamelCase(member.Name), GetTypeDefaultValue(member))
			reqFeilds = append(reqFeilds, str)
		}
	}
	reqFeildsStr := strings.Join(reqFeilds, "\n")
	if len(reqFeilds) > 0 {
		reqFeildsStr = "\n" + reqFeildsStr + "\n"
	}
	reqFeildsStr = "{" + reqFeildsStr + "}"

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          subDir,
		filename:        goFile + "_test.go",
		templateName:    "logicTemplate",
		category:        category,
		templateFile:    logicTemplateTestFile,
		builtinTemplate: logicTestTemplate,
		data: map[string]any{
			"pkgName":        subDir[strings.LastIndex(subDir, "/")+1:],
			"ProjectName":    filepath.Base(rootPkg),
			"NoTypes":        !strings.Contains(importPackages, "/types"),
			"ImportPackages": importPackages,
			"imports":        imports,
			"logic":          strings.Title(logic),
			"LogicType":      strings.Title(getLogicName(route)),
			"function":       strings.Title(strings.TrimSuffix(logic, "Logic")),
			"LogicName":      logicName,
			"responseType":   responseString,
			"returnString":   returnString,
			"request":        requestString,
			"Call":           strings.Title(strings.TrimSuffix(handler, "Handler")),
			"RequestType":    util.Title(route.RequestTypeName()),
			"ResponseType":   util.Title(responseGoTypeName(route, typesPacket)),
			"hasDoc":         len(route.JoinedDoc()) > 0,
			"HasResp":        len(route.ResponseTypeName()) > 0,
			"HasRequest":     len(route.RequestTypeName()) > 0,
			"doc":            getDoc(route.JoinedDoc()),
			"config":         fmt.Sprintf("\"%s\"", pathx.JoinPackages(rootPkg, configDir)),
			"types":          fmt.Sprintf("\"%s\"", pathx.JoinPackages(rootPkg, typesDir)),
			"reqFeilds":      reqFeildsStr,
		},
	})
}

func getLogicFolderPathTest(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return logicDirTest
		}
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(logicDirTest, folder)
}

func GetTypeDefaultValue(member spec.Member) string {
	switch name := member.Type.Name(); name {
	// 整型及别名（包括有符号、无符号、指针类型）
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune":
		if strings.Contains(strings.ToLower(member.Name), "size") {
			return "10"
		}
		if strings.Contains(strings.ToLower(member.Name), "page") {
			return "1"
		}
		return "0"
	// 浮点型
	case "float32", "float64":
		return "0.0"
	// 字符串
	case "string":
		return "\"\""
	// 布尔值
	case "bool":
		return "false"
	// 复数类型
	case "complex64", "complex128":
		return fmt.Sprintf("%v", complex(0, 0))
	// 引用类型及指针
	case "slice", "map", "chan", "func", "pointer", "interface":
		return "nil"
	// 空结构体默认值
	case "struct":
		return "struct{}{}"
	default:
		if strings.HasPrefix(name, "*") {
			return fmt.Sprintf("%s{}", strings.Replace(name, "*", "&types.", 1))
		}
		if strings.HasPrefix(name, "[]") {
			if strings.Contains(name, "*") {
				return fmt.Sprintf("%s{}", strings.Replace(name, "*", "*types.", 1))
			} else {
				switch name {
				// 整型及别名（包括有符号、无符号、指针类型）
				case "int", "int8", "int16", "int32", "int64",
					"uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune":
					fallthrough
				// 浮点型
				case "float32", "float64":
					fallthrough
				// 字符串
				case "string":
					fallthrough
				// 布尔值
				case "bool":
					fallthrough
				// 复数类型
				case "complex64", "complex128":
					fallthrough
				// 引用类型及指针
				case "slice", "map", "chan", "func", "pointer", "interface":
					fallthrough
				// 空结构体默认值
				case "struct":
					return fmt.Sprintf("%s{}", name)
				default:
					return fmt.Sprintf("[]types.%s{}", strings.TrimPrefix(name, "[]"))
				}
			}
		}
		return fmt.Sprintf("types.%s{}", name)
	}
}

// 大驼峰命名法 首字母大写(暂时)
func upperCamelCase(word string) string {
	if len(word) == 0 {
		return ""
	}
	return strings.ToUpper(string(word[0])) + word[1:]
}
