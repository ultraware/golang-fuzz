package main

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

func generateLibFuzzer(sourceCode *ast.File, fFuzz *ast.FuncLit) {
	fuzzFile, err := os.CreateTemp(`.`, `fuzz-*_test.go`)
	if err != nil {
		panic(err)
	}
	defer func() {
		fuzzFile.Close()
		if !*keepFile {
			os.Remove(fuzzFile.Name())
		}
	}()

	// Add build tag to prevent conflicts
	if len(sourceCode.Comments) == 0 {
		sourceCode.Comments[0] = &ast.CommentGroup{}
	}
	sourceCode.Comments[0].List = append(sourceCode.Comments[0].List, nil)
	copy(sourceCode.Comments[0].List[1:], sourceCode.Comments[0].List)
	sourceCode.Comments[0].List[0] = &ast.Comment{Text: `//go:build gofuzzbuild`}

	funcDecl := &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "GoFuzzBuildFuzz",
		},
		Type: fFuzz.Type,
		Body: fFuzz.Body,
	}

	// Add Fuzz function
	sourceCode.Decls = append(sourceCode.Decls, funcDecl)

	// Write to file
	err = format.Node(fuzzFile, token.NewFileSet(), sourceCode)
	if err != nil {
		panic(err)
	}
}

// TODO: t *testing.T -> New struct in top (var t = ...)
