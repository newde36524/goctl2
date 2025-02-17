package generator

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	conf "github.com/newde36524/goctl2/config"
	"github.com/newde36524/goctl2/rpc/parser"
	"github.com/newde36524/goctl2/util"
	"github.com/newde36524/goctl2/util/format"
	"github.com/newde36524/goctl2/util/pathx"
	"github.com/newde36524/goctl2/util/stringx"
	"github.com/zeromicro/go-zero/core/collection"
)

//go:embed logic.test.tpl
var logicTestTemplate string

// GenLogic generates the logic file of the rpc service, which corresponds to the RPC definition items in proto.
func (g *Generator) GenLogicTest(ctx DirContext, proto parser.Proto, cfg *conf.Config,
	c *ZRpcContext) error {
	if !c.Multiple {
		return g.genLogicInCompatibility2(ctx, proto, cfg)
	}

	return g.genLogicGroup2(ctx, proto, cfg)
}

func (g *Generator) genLogicInCompatibility2(ctx DirContext, proto parser.Proto,
	cfg *conf.Config) error {
	dir := ctx.GetLogicTest()
	service := proto.Service[0].Service.Name
	for _, rpc := range proto.Service[0].RPC {
		logicName := fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel())
		logicFilename, err := format.FileNamingFormat(cfg.NamingFormat, rpc.Name+"_logic")
		if err != nil {
			return err
		}

		filename := filepath.Join(dir.Filename, logicFilename+".go")
		functions, err := g.genLogicFunction2(service, proto.PbPackage, logicName, rpc)
		if err != nil {
			return err
		}

		imports := collection.NewSet()
		imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
		imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetPb().Package))
		text, err := pathx.LoadTemplate(category, logicTemplateFileFile, logicTemplate)
		if err != nil {
			return err
		}
		err = util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]any{
			"logicName":   fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel()),
			"functions":   functions,
			"packageName": "logic",
			"imports":     strings.Join(imports.KeysStr(), pathx.NL),
		}, filename, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) genLogicGroup2(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetLogicTest()
	for _, item := range proto.Service {
		serviceName := item.Name
		for _, rpc := range item.RPC {
			var (
				err           error
				filename      string
				logicName     string
				logicFilename string
				packageName   string
			)

			logicName = fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel())
			childPkg, err := dir.GetChildPackage(serviceName)
			if err != nil {
				return err
			}

			serviceDir := filepath.Base(childPkg)
			nameJoin := fmt.Sprintf("%s_logic", serviceName)
			packageName = strings.ToLower(stringx.From(nameJoin).ToCamel())
			logicFilename, err = format.FileNamingFormat(cfg.NamingFormat, rpc.Name+"_logic")
			if err != nil {
				return err
			}

			filename = filepath.Join(dir.Filename, serviceDir, logicFilename+".go")
			functions, err := g.genLogicFunction2(serviceName, proto.PbPackage, logicName, rpc)
			if err != nil {
				return err
			}

			imports := collection.NewSet()
			imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
			imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetPb().Package))
			text, err := pathx.LoadTemplate(category, logicTestTemplateFileFile, logicTestTemplate)
			if err != nil {
				return err
			}

			if err = util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]any{
				"ProjectName":  proto.PbPackage,
				"logicName":    logicName,
				"functions":    functions,
				"serviceName":  strings.ToLower(stringx.From(serviceName).ToCamel()),
				"method":       parser.CamelCase(rpc.Name),
				"packageName":  packageName,
				"imports":      strings.Join(imports.KeysStr(), pathx.NL),
				"hasReq":       !rpc.StreamsRequest,
				"hasReply":     !rpc.StreamsRequest && !rpc.StreamsReturns,
				"stream":       rpc.StreamsRequest || rpc.StreamsReturns,
				"request":      fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(rpc.RequestType)),
				"response":     fmt.Sprintf("*%s.%s", proto.PbPackage, parser.CamelCase(rpc.ReturnsType)),
				"responseType": fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(rpc.ReturnsType)),
			}, filename, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) genLogicFunction2(serviceName, goPackage, logicName string,
	rpc *parser.RPC) (string,
	error) {
	functions := make([]string, 0)
	text, err := pathx.LoadTemplate(category, logicFuncTemplateFileFile, logicTestTemplate)
	if err != nil {
		return "", err
	}

	comment := parser.GetComment(rpc.Doc())
	streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(serviceName),
		parser.CamelCase(rpc.Name), "Server")
	buffer, err := util.With("fun").Parse(text).Execute(map[string]any{
		"logicName":    logicName,
		"method":       parser.CamelCase(rpc.Name),
		"hasReq":       !rpc.StreamsRequest,
		"request":      fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.RequestType)),
		"hasReply":     !rpc.StreamsRequest && !rpc.StreamsReturns,
		"response":     fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
		"responseType": fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
		"stream":       rpc.StreamsRequest || rpc.StreamsReturns,
		"streamBody":   streamServer,
		"hasComment":   len(comment) > 0,
		"comment":      comment,
	})
	if err != nil {
		return "", err
	}

	functions = append(functions, buffer.String())
	return strings.Join(functions, pathx.NL), nil
}
