package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gookit/color"
)

var ErrMissingService = errors.New("missing service")

// Validate validates Validate the integrity of the spec.
func (s *ApiSpec) Validate() error {
	if len(s.Service.Name) == 0 {
		return ErrMissingService
	}
	if len(s.Service.Groups) == 0 {
		return ErrMissingService
	}
	return nil
}

func (s *ApiSpec) tagClear(tag string) string {
	vTag := strings.Trim(tag, "`")
	vTag = strings.Trim(vTag, `"`)
	vTag = strings.TrimPrefix(vTag, "form:")
	vTag = strings.TrimPrefix(vTag, "json:")
	vTag = strings.ReplaceAll(vTag, "\"", "")
	vTag = strings.ReplaceAll(vTag, `"`, "")
	return vTag
}

// CheckTag 检查tag是否合法
func (s *ApiSpec) CheckTag() {
	var (
		gets      []string
		posts     []string
		tagFatail []string
		mp        = map[string]string{
			"get":  "form",
			"post": "json",
		}
	)
	for _, group := range s.Service.Groups {
		for _, route := range group.Routes {
			method := strings.ToLower(route.Method)
			if tag, ok := mp[method]; ok {
				for _, _type := range s.Types {
					if route.RequestType == nil {
						continue
					}
					if _type.Name() == route.RequestType.Name() {
						v := _type.(DefineStruct)
						for _, member := range v.Members {
							memberName := member.Name
							if len(memberName) == 0 { //直接组合的类型不用判断tag
								continue
							}
							vTag := s.tagClear(member.Tag)
							if len(vTag) == 0 {
								tagFatail = append(tagFatail, color.Red.Render(fmt.Errorf("tag is not set. [member: %s] [typeName: %s] [group: %s] [prefix: %s] [tags: %s]", memberName, _type.Name(), group.Annotation.Properties["group"], group.Annotation.Properties["prefix"], group.Annotation.Properties["tags"]).Error()))
								continue
							}
							firstChar := strings.Split(vTag, "")[0]
							if firstChar == strings.ToUpper(firstChar) { //首字母必须小写
								tagFatail = append(tagFatail, color.Red.Render(fmt.Errorf("tag first char is Upper, now have %s. [member: %s] [typeName: %s] [group: %s] [prefix: %s] [tags: %s]", member.Tag, memberName, _type.Name(), group.Annotation.Properties["group"], group.Annotation.Properties["prefix"], group.Annotation.Properties["tags"]).Error()))
								continue
							}
							if !strings.Contains(member.Tag, tag) {
								if method == "get" {
									gets = append(gets, color.Red.Render(fmt.Errorf("%s method must use `%s` tag for %s, now have %s. [member: %s] [typeName: %s] [group: %s] [prefix: %s] [tags: %s]", method, tag, route.Path, member.Tag, memberName, _type.Name(), group.Annotation.Properties["group"], group.Annotation.Properties["prefix"], group.Annotation.Properties["tags"]).Error()))
								}
								if method == "post" {
									posts = append(posts, color.Yellow.Render(fmt.Errorf("%s method must use `%s` tag for %s, now have %s. [member: %s] [typeName: %s] [group: %s] [prefix: %s] [tags: %s]", method, tag, route.Path, member.Tag, memberName, _type.Name(), group.Annotation.Properties["group"], group.Annotation.Properties["prefix"], group.Annotation.Properties["tags"]).Error()))
								}
							}
						}
					}
				}
			}
		}
	}
	s.onceCheck(append(append(posts, tagFatail...), gets...), func(s string) {
		fmt.Println(s)
	})
	// for _, v := range gets {
	// 	fmt.Println(v)
	// }
}

// onceCheck 特殊处理 旧的接口只提示一次，往后新增的接口再提示
func (s *ApiSpec) onceCheck(ignoreList []string, fn func(s string)) {
	fileName := "ignoreTagList"
	ex, _ := os.Executable()
	fullName := filepath.Join(filepath.Dir(ex), fileName)
	if _, err := os.Stat(fullName); os.IsNotExist(err) {
		fs, _ := os.Create(fullName)
		bs, _ := json.Marshal(ignoreList)
		fs.WriteString(string(bs))
		fs.Close()
		for _, v := range ignoreList {
			fn(v)
		}
		return
	}
	bs, _ := os.ReadFile(fullName)
	var ignoreList2 []string
	json.Unmarshal(bs, &ignoreList2)
	tmp := false
	for _, v := range ignoreList {
		if slices.Contains(ignoreList2, v) {
			tmp = true
			continue
		}
		fn(v)
	}
	if tmp {
		fmt.Println("已忽略旧不规范文档，如需再次检测可执行命令删除记录文件 rm -rf $(go env GOPATH)/bin/" + fileName)
	}
}
