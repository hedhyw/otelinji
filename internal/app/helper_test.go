package app

import (
	"testing"

	"github.com/dave/dst"
	"github.com/stretchr/testify/assert"
)

func TestFuncName(t *testing.T) {
	t.Parallel()

	assert.Empty(t, funcName(nil), "nil_func")
	assert.Empty(t, funcName(&dst.FuncDecl{Name: nil}), "nil_name")
	assert.Equal(t, "name", funcName(&dst.FuncDecl{Name: dst.NewIdent("name")}), "ok")
}

func TestReceiverType(t *testing.T) {
	t.Parallel()

	assert.Empty(t, receiverType(nil), "nil_func")
	assert.Empty(t, receiverType(&dst.FuncDecl{Recv: nil}), "nil_func")
	assert.Empty(t, receiverType(&dst.FuncDecl{
		Recv: &dst.FieldList{List: []*dst.Field{{Type: nil}}}},
	), "nil_type")
	assert.Empty(t, receiverType(&dst.FuncDecl{
		Recv: &dst.FieldList{List: []*dst.Field{}}},
	), "empty_list")
	assert.Equal(t, "int", receiverType(&dst.FuncDecl{
		Recv: &dst.FieldList{List: []*dst.Field{{Type: dst.NewIdent("int")}}}},
	), "ok")
}

func TestContextParamNameFromFunc(t *testing.T) {
	t.Parallel()

	t.Run("failed", func(t *testing.T) {
		t.Parallel()

		_, ok := contextParamNameFromFunc(nil)
		assert.False(t, ok, "nil_func")

		_, ok = contextParamNameFromFunc(&dst.FuncDecl{Type: nil})
		assert.False(t, ok, "nil_type")

		_, ok = contextParamNameFromFunc(&dst.FuncDecl{Type: &dst.FuncType{Params: nil}})
		assert.False(t, ok, "nil_params")

		_, ok = contextParamNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Params: &dst.FieldList{List: []*dst.Field{nil}}}},
		)
		assert.False(t, ok, "nil_field")

		_, ok = contextParamNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Params: &dst.FieldList{List: []*dst.Field{{
				Type: dst.NewIdent("int"),
			}}}},
		})
		assert.False(t, ok, "int_type")
	})

	t.Run("ok_named", func(t *testing.T) {
		t.Parallel()

		name, ok := contextParamNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Params: &dst.FieldList{List: []*dst.Field{{
				Type:  dst.NewIdent("context.Context"),
				Names: []*dst.Ident{dst.NewIdent("ctxTest")},
			}}}}},
		)
		if assert.True(t, ok) {
			assert.Equal(t, "ctxTest", name)
		}
	})

	t.Run("ok_unnamed", func(t *testing.T) {
		t.Parallel()

		name, ok := contextParamNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Params: &dst.FieldList{List: []*dst.Field{{
				Type: dst.NewIdent("context.Context"),
			}}}}},
		)
		if assert.True(t, ok) {
			assert.Equal(t, "ctx", name)
		}
	})
}

func TestFieldType(t *testing.T) {
	t.Parallel()

	assert.Empty(t, fieldType(nil), "nil_field")
	assert.Empty(t, fieldType(&dst.Field{Type: nil}), "nil_type")
	assert.Equal(t, "context.Context", fieldType(&dst.Field{
		Type: &dst.SelectorExpr{
			X:   dst.NewIdent("context"),
			Sel: dst.NewIdent("Context"),
		},
	}), "ok_selector")
	assert.Equal(t, "int", fieldType(&dst.Field{
		Type: dst.NewIdent("int"),
	}), "ok_type")
}

func TestErrResultNameFromFunc(t *testing.T) {
	t.Parallel()

	t.Run("failed_nil", func(t *testing.T) {
		t.Parallel()

		assert.Empty(t, errResultNameFromFunc(nil), "nil_func")
		assert.Empty(t, errResultNameFromFunc(&dst.FuncDecl{Type: nil}), "nil_type")
		assert.Empty(t, errResultNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Results: nil},
		}), "nil_results")
		assert.Empty(t, errResultNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Results: &dst.FieldList{List: nil}},
		}), "nil_;ist")
		assert.Empty(t, errResultNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Results: &dst.FieldList{List: []*dst.Field{nil}}},
		}), "nil_result_item")

	})

	t.Run("failed_invalid", func(t *testing.T) {
		t.Parallel()

		assert.Empty(t, errResultNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Results: &dst.FieldList{List: []*dst.Field{{
				Type: &dst.BadExpr{},
			}}}},
		}), "not_an_indent")
		assert.Empty(t, errResultNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Results: &dst.FieldList{List: []*dst.Field{{
				Type:  dst.NewIdent("int"),
				Names: []*dst.Ident{dst.NewIdent("err")},
			}}}},
		}), "not_an_error")
		assert.Empty(t, errResultNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Results: &dst.FieldList{List: []*dst.Field{{
				Type: &dst.Ident{
					Path: "path",
				},
				Names: []*dst.Ident{dst.NewIdent("err")},
			}}}},
		}), "has_path")
		assert.Empty(t, errResultNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Results: &dst.FieldList{List: []*dst.Field{{
				Type:  dst.NewIdent("error"),
				Names: []*dst.Ident{dst.NewIdent("err1"), dst.NewIdent("err2")},
			}}}},
		}), "many_names")
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "err", errResultNameFromFunc(&dst.FuncDecl{
			Type: &dst.FuncType{Results: &dst.FieldList{List: []*dst.Field{{
				Type:  dst.NewIdent("error"),
				Names: []*dst.Ident{dst.NewIdent("err")},
			}}}},
		}), "many_names")
	})
}

func TestAddImportSpecs(t *testing.T) {
	t.Parallel()

	t.Run("nil_file", func(t *testing.T) {
		t.Parallel()

		addImportSpecs(nil, []*dst.ImportSpec{{}})
	})

	t.Run("empty_spec", func(t *testing.T) {
		t.Parallel()

		addImportSpecs(&dst.File{}, nil)
	})

	t.Run("ok_new", func(t *testing.T) {
		t.Parallel()

		f := &dst.File{}

		addImportSpecs(f, []*dst.ImportSpec{{
			Name: dst.NewIdent("fmt"),
			Path: &dst.BasicLit{
				Value: "fmt",
			},
		}})

		if assert.Len(t, f.Imports, 1) {
			assert.Equal(t, "fmt", f.Imports[0].Name.Name)
		}
	})

	t.Run("ok_found", func(t *testing.T) {
		t.Parallel()

		f := &dst.File{
			Imports: []*dst.ImportSpec{{
				Name: dst.NewIdent("fmt"),
				Path: &dst.BasicLit{
					Value: "fmt",
				},
			}},
		}

		addImportSpecs(f, []*dst.ImportSpec{{
			Name: dst.NewIdent("fmt"),
			Path: &dst.BasicLit{
				Value: "fmt",
			},
		}})

		if assert.Len(t, f.Imports, 1) {
			assert.Equal(t, "fmt", f.Imports[0].Name.Name)
		}
	})
}

func TestFilterUsedPackages(t *testing.T) {
	t.Parallel()

	t.Run("nil_file", func(t *testing.T) {
		t.Parallel()

		resSpec, resImport := filterUsedPackages(nil, []*dst.ImportSpec{{}})
		assert.Nil(t, resSpec, "resSpec")
		assert.Nil(t, resImport, "resImport")
	})

	t.Run("empty_spec", func(t *testing.T) {
		t.Parallel()

		resSpec, resImport := filterUsedPackages(&dst.File{}, nil)
		assert.Nil(t, resSpec, "resSpec")
		assert.Nil(t, resImport, "resImport")
	})

	t.Run("ok_new", func(t *testing.T) {
		t.Parallel()

		resSpec, resImport := filterUsedPackages(&dst.File{}, []*dst.ImportSpec{{
			Name: dst.NewIdent("fmt"),
			Path: &dst.BasicLit{Value: "fmt"},
		}})
		assert.Len(t, resSpec, 1, "resSpec")
		assert.Len(t, resImport, 1, "resImport")
	})

	t.Run("ok_filtered", func(t *testing.T) {
		t.Parallel()

		resSpec, resImport := filterUsedPackages(&dst.File{
			Imports: []*dst.ImportSpec{{
				Name: dst.NewIdent("fmt"),
				Path: &dst.BasicLit{Value: "fmt"},
			}},
		}, []*dst.ImportSpec{{
			Name: dst.NewIdent("fmt"),
			Path: &dst.BasicLit{Value: "fmt"},
		}})
		assert.Empty(t, resSpec, "resSpec")
		assert.Empty(t, resImport, "resImport")
	})
}

func TestIsCImportSpec(t *testing.T) {
	t.Parallel()

	assert.False(t, isCImportSpec(nil), "nil")
	assert.False(t, isCImportSpec(&dst.TypeSpec{}), "not_at_import_spec")
	assert.False(t, isCImportSpec((*dst.ImportSpec)(nil)), "nil_import_spec")
	assert.False(t, isCImportSpec(&dst.ImportSpec{Path: nil}), "nil_path")
	assert.False(t, isCImportSpec(&dst.ImportSpec{
		Path: &dst.BasicLit{Value: "fmt"},
	}), "not_a_c")
	assert.False(t, isCImportSpec(&dst.ImportSpec{
		Path: &dst.BasicLit{Value: "C"},
	}), "no_quotes")
	assert.False(t, isCImportSpec(&dst.ImportSpec{
		Path: &dst.BasicLit{Value: `"c"`},
	}), "small_c")
	assert.True(t, isCImportSpec(&dst.ImportSpec{
		Path: &dst.BasicLit{Value: `"C"`},
	}), "ok_c")
}

func TestIsGenerated(t *testing.T) {
	t.Parallel()

	assert.False(t, isGenerated(nil), "nil_node")
	assert.False(t, isGenerated(&dst.BasicLit{}), "empty")
	assert.False(t, isGenerated(&dst.BlockStmt{List: []dst.Stmt{}}), "nil_item")
	assert.True(t, isGenerated(&dst.BlockStmt{List: []dst.Stmt{
		&dst.DeclStmt{Decs: dst.DeclStmtDecorations{
			NodeDecs: dst.NodeDecs{Start: dst.Decorations{"DO NOT EDIT"}},
		}}},
	}), "generated_start")
	assert.True(t, isGenerated(&dst.BlockStmt{List: []dst.Stmt{
		&dst.DeclStmt{Decs: dst.DeclStmtDecorations{
			NodeDecs: dst.NodeDecs{End: dst.Decorations{"DO NOT EDIT"}},
		}}},
	}), "generated_end")
}
