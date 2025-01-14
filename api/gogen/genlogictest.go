package gogen

import (
	_ "embed"
	"path"
	"path/filepath"
	"strings"

	"github.com/newde36524/goctl2/api/spec"
	"github.com/newde36524/goctl2/config"
	"github.com/newde36524/goctl2/util"
	"github.com/newde36524/goctl2/util/format"
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
