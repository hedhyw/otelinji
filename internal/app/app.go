package app

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"text/template"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/hedhyw/otelinji/internal/pkg/assets"
	"github.com/hedhyw/otelinji/internal/pkg/config"
)

// App command-line.
type App struct {
	cfg        *config.Config
	fset       *token.FileSet
	injectTmpl *template.Template
}

// New app.
func New(cfg *config.Config) *App {
	return &App{
		cfg:  cfg,
		fset: token.NewFileSet(),
	}
}

// Run the application and print the result to out.
func (a *App) Run(out io.Writer) (err error) {
	err = a.preapreInjectTemplate()
	if err != nil {
		return fmt.Errorf("reading inject template: %w", err)
	}

	inputData, err := os.ReadFile(a.cfg.FileName)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	dstFile, err := a.parseFile(inputData)
	if err != nil {
		return fmt.Errorf("parsing file: %w", err)
	}

	if a.cfg.SkipGenerated && isGenerated(dstFile) {
		if a.cfg.WriteIntoFile {
			return nil
		}

		_, err = out.Write(inputData)

		return err
	}

	if err = a.processFile(dstFile); err != nil {
		return fmt.Errorf("processing file: %w", err)
	}

	if err = a.printNode(out, dstFile); err != nil {
		return fmt.Errorf("printing ast file: %w", err)
	}

	return nil
}

func (a *App) preapreInjectTemplate() (err error) {
	tmplContent, err := getTemplateContent(a.cfg.TemplateName)
	if err != nil {
		return fmt.Errorf("getting template content: %w", err)
	}

	tmpl, err := template.New("template_source").
		Funcs(getTemplateFunc()).
		Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("parsing inject template: %w", err)
	}

	a.injectTmpl = tmpl

	return nil
}

func (a *App) parseFile(data []byte) (fstFile *dst.File, err error) {
	astFile, err := parser.ParseFile(
		a.fset,
		a.cfg.FileName,
		data,
		parser.ParseComments,
	)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}

	dec := decorator.NewDecorator(a.fset)

	dstFile, err := dec.DecorateFile(astFile)
	if err != nil {
		return nil, fmt.Errorf("decorate file: %w", err)
	}

	return dstFile, nil
}

func (a *App) processFile(dstFile *dst.File) (err error) {
	for _, decl := range dstFile.Decls {
		if fnDecl, ok := decl.(*dst.FuncDecl); ok {
			fnName := funcName(fnDecl)

			if !dst.IsExported(fnName) {
				continue
			}

			if err = a.injectBlock(dstFile, fnDecl); err != nil {
				return fmt.Errorf("injecting block: %s: %w", fnName, err)
			}
		}
	}

	return nil
}

func (a *App) printNode(out io.Writer, dstFile *dst.File) (err error) {
	if a.cfg.WriteIntoFile {
		//nolint: nosnakecase // Std flag.
		f, err := os.OpenFile(a.cfg.FileName, os.O_WRONLY, os.ModePerm)
		if err != nil {
			return fmt.Errorf("opening file to write: %w", err)
		}

		defer func() { err = semerr.NewMultiError(err, f.Close()) }()

		out = f
	}

	var astBuf bytes.Buffer

	restorer := decorator.NewRestorer()

	if err = restorer.Fprint(&astBuf, dstFile); err != nil {
		return fmt.Errorf("fprint: %w", err)
	}

	fmtSource, err := format.Source(astBuf.Bytes())
	if err != nil {
		return fmt.Errorf("format source: %w", err)
	}

	if _, err = out.Write(fmtSource); err != nil {
		return fmt.Errorf("write source: %w", err)
	}

	return nil
}

func (a *App) injectBlock(
	dstFile *dst.File,
	fnDecl *dst.FuncDecl,
) (err error) {
	ctxParamName, _ := contextParamNameFromFunc(fnDecl)
	fnName := funcName(fnDecl)

	injStms, err := a.getInjectClause(map[string]any{
		// Keep it in sync with the documentation in the template.
		"CtxParamName":  ctxParamName,
		"FuncName":      fnName,
		"PackageName":   dstFile.Name.Name,
		"ReceiverType":  receiverType(fnDecl),
		"IsContextUsed": checkCtxUsed(fnDecl.Body, ctxParamName),
		"ErrResultName": errResultNameFromFunc(fnDecl),
	})
	if err != nil {
		return fmt.Errorf("getting inj clouse: %w", err)
	}

	if checkIndents(fnDecl, injStms.Condition) {
		return nil
	}

	if len(injStms.Body) > 0 {
		addImportSpecs(dstFile, injStms.ImportSpecs)
	}

	fnDecl.Body.List = append(injStms.Body, fnDecl.Body.List...)

	return nil
}

func getTemplateContent(fileName string) (string, error) {
	if fileName == "@/otel" || fileName == "" {
		return assets.OtelTmpl(), nil
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	return string(data), nil
}
