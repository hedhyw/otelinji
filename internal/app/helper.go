package app

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/dave/dst"
)

func funcName(fnDecl *dst.FuncDecl) string {
	if fnDecl == nil || fnDecl.Name == nil {
		return ""
	}

	return fnDecl.Name.String()
}

func receiverType(fnDecl *dst.FuncDecl) string {
	if fnDecl == nil || fnDecl.Recv == nil {
		return ""
	}

	for _, r := range fnDecl.Recv.List {
		switch t := r.Type.(type) {
		case nil:
			return ""
		case *dst.StarExpr:
			return fmt.Sprint(t.X)
		default:
			return fmt.Sprint(t)
		}
	}

	return ""
}

func contextParamNameFromFunc(fnDecl *dst.FuncDecl) (ctxName string, ok bool) {
	const (
		typeContext = "context.Context"
		defContext  = "ctx"
	)

	if fnDecl == nil || fnDecl.Type == nil || fnDecl.Type.Params == nil {
		return "", false
	}

	for _, param := range fnDecl.Type.Params.List {
		if param == nil {
			continue
		}

		paramType := fieldType(param)
		if paramType != typeContext {
			continue
		}

		if len(param.Names) == 0 || param.Names[0] == nil {
			convertUnnamedToDashes(fnDecl.Type.Params.List)

			param.Names = []*dst.Ident{
				dst.NewIdent(defContext),
			}

			return defContext, true
		}

		paramName := param.Names[0].String()

		if paramName == "_" {
			paramName = defContext
			param.Names[0].Name = paramName
		}

		return paramName, true
	}

	return "", false
}

func convertUnnamedToDashes(fields []*dst.Field) {
	for _, param := range fields {
		if param == nil {
			continue
		}

		if len(param.Names) == 0 || param.Names[0] == nil {
			param.Names = []*dst.Ident{
				dst.NewIdent("_"),
			}
		}
	}
}

func fieldType(param *dst.Field) string {
	if param == nil || param.Type == nil {
		return ""
	}

	if selExpr, ok := param.Type.(*dst.SelectorExpr); ok {
		return strings.Join([]string{
			fmt.Sprint(selExpr.X),
			fmt.Sprint(selExpr.Sel),
		}, ".")
	}

	return fmt.Sprint(param.Type)
}

func checkCtxUsed(node dst.Node, paramName string) (used bool) {
	dst.Inspect(node, func(n dst.Node) bool {
		if nodeIndent, ok := n.(*dst.Ident); ok {
			if nodeIndent.Name == paramName {
				used = true

				return false
			}
		}

		return true
	})

	return used
}

func errResultNameFromFunc(fnDecl *dst.FuncDecl) (errParamName string) {
	if fnDecl == nil || fnDecl.Type == nil || fnDecl.Type.Results == nil {
		return ""
	}

	for _, result := range fnDecl.Type.Results.List {
		if result == nil {
			continue
		}

		indentType, ok := result.Type.(*dst.Ident)
		if !ok {
			continue
		}

		if indentType.Name != "error" || indentType.Path != "" || len(result.Names) != 1 {
			continue
		}

		if result.Names[0].Name == "_" {
			return ""
		}

		return result.Names[0].Name
	}

	return ""
}

func addImportSpecs(file *dst.File, specs []*dst.ImportSpec) {
	if len(specs) == 0 || file == nil {
		return
	}

	targetSpecs, targetImportSpecs := filterUsedPackages(file, specs)

	if len(targetSpecs) == 0 {
		return
	}

	var added bool

	dst.Inspect(file, func(n dst.Node) bool {
		if genDecl, ok := n.(*dst.GenDecl); ok && genDecl.Tok == token.IMPORT {
			if len(genDecl.Specs) == 1 && isCImportSpec(genDecl.Specs[0]) {
				return true
			}

			genDecl.Specs = append(genDecl.Specs, targetSpecs...)
			file.Imports = append(file.Imports, targetImportSpecs...)
			added = true

			return false
		}

		return true
	})

	if !added {
		file.Decls = append([]dst.Decl{&dst.GenDecl{
			Tok:   token.IMPORT,
			Specs: targetSpecs,
		}}, file.Decls...)
		file.Imports = append(file.Imports, targetImportSpecs...)
	}
}

func filterUsedPackages(file *dst.File, specs []*dst.ImportSpec) ([]dst.Spec, []*dst.ImportSpec) {
	if len(specs) == 0 || file == nil {
		return nil, nil
	}

	specsMap := make(map[string]*dst.ImportSpec, len(specs))

	for _, s := range specs {
		specsMap[s.Path.Value] = s
	}

	for _, pkg := range file.Imports {
		if pkg != nil && pkg.Path != nil {
			delete(specsMap, pkg.Path.Value)
		}
	}

	if len(specsMap) == 0 {
		return nil, nil
	}

	targetImportSpecs := make([]*dst.ImportSpec, 0, len(specsMap))
	targetSpecs := make([]dst.Spec, 0, len(specsMap))

	for _, s := range specsMap {
		targetSpecs = append(targetSpecs, s)
		targetImportSpecs = append(targetImportSpecs, s)
	}

	return targetSpecs, targetImportSpecs
}

func isCImportSpec(spec dst.Spec) bool {
	const importPathC = `"C"`

	impSpec, ok := spec.(*dst.ImportSpec)
	if !ok {
		return false
	}

	if impSpec == nil || impSpec.Path == nil {
		return false
	}

	return impSpec.Path.Value == importPathC
}

func isGenerated(f dst.Node) (isGen bool) {
	const commentGenerated = "DO NOT EDIT"

	if f == nil {
		return false
	}

	dst.Inspect(f, func(n dst.Node) bool {
		if n == nil {
			return true
		}

		decorations := n.Decorations()
		if decorations == nil {
			return true
		}

		if decorations.End != nil {
			if anyContains(decorations.End.All(), commentGenerated) {
				isGen = true

				return false
			}
		}

		if decorations.Start != nil {
			if anyContains(decorations.Start.All(), commentGenerated) {
				isGen = true

				return false
			}
		}

		return true
	})

	return isGen
}

func anyContains(val []string, substr string) bool {
	for _, v := range val {
		if strings.Contains(v, substr) {
			return true
		}
	}

	return false
}
