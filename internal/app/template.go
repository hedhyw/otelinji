package app

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"
	"text/template"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

type injectClause struct {
	Body        []dst.Stmt
	Condition   []string
	ImportSpecs []*dst.ImportSpec
}

const (
	errNilFuncDecl semerr.Error = "nil func decl"
)

func (a *App) getInjectClause(params map[string]any) (*injectClause, error) {
	var renderBuf bytes.Buffer

	if err := a.injectTmpl.Execute(&renderBuf, params); err != nil {
		return nil, fmt.Errorf("execing template: %w", err)
	}

	dstFile, err := decorator.Parse(renderBuf.String())
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	var fnDecl *dst.FuncDecl

	dst.Inspect(dstFile, func(n dst.Node) bool {
		if fd, ok := n.(*dst.FuncDecl); ok {
			fnDecl = fd

			return false
		}

		return true
	})

	if fnDecl == nil || fnDecl.Body == nil {
		return nil, errNilFuncDecl
	}

	return &injectClause{
		Body:        fnDecl.Body.List,
		Condition:   getCheckIndentsCondition(dstFile),
		ImportSpecs: extractImportSpecs(dstFile),
	}, nil
}

func extractImportSpecs(node dst.Node) []*dst.ImportSpec {
	specs := []dst.Spec{}

	dst.Inspect(node, func(n dst.Node) bool {
		if genDecl, ok := n.(*dst.GenDecl); ok && genDecl.Tok == token.IMPORT {
			specs = append(specs, genDecl.Specs...)
		}

		return true
	})

	if len(specs) == 0 {
		return nil
	}

	importSpecs := make([]*dst.ImportSpec, 0, len(specs))

	for _, s := range specs {
		if impSpec, ok := s.(*dst.ImportSpec); ok {
			importSpecs = append(importSpecs, impSpec)
		}
	}

	return importSpecs
}

func getCheckIndentsCondition(node dst.Node) []string {
	const prefixCheck = "//otelinji:check-indents "

	var comment string

	dst.Inspect(node, func(n dst.Node) bool {
		if n == nil {
			return false
		}

		for _, nodeComment := range n.Decorations().Start.All() {
			if strings.HasPrefix(nodeComment, prefixCheck) {
				comment = strings.TrimPrefix(nodeComment, prefixCheck)

				return false
			}
		}

		return true
	})

	if comment == "" {
		return nil
	}

	return strings.Split(comment, ",")
}

func checkIndents(node dst.Node, condition []string) (found bool) {
	if len(condition) == 0 {
		return false
	}

	conditionSet := make(map[string]struct{}, len(condition))

	for _, c := range condition {
		conditionSet[c] = struct{}{}
	}

	dst.Inspect(node, func(n dst.Node) bool {
		if nodeIndent, ok := n.(*dst.Ident); ok {
			if _, ok = conditionSet[nodeIndent.Name]; ok {
				delete(conditionSet, nodeIndent.Name)
			}

			return len(conditionSet) > 0
		}

		return true
	})

	return len(conditionSet) == 0
}

func getTemplateFunc() template.FuncMap {
	return template.FuncMap{
		"joinWithDot": func(elems ...string) string {
			nonBlank := make([]string, 0, len(elems))

			for _, el := range elems {
				if el != "" {
					nonBlank = append(nonBlank, el)
				}
			}

			return strings.Join(nonBlank, ".")
		},
	}
}
