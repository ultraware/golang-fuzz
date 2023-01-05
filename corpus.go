package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	_ "embed"
)

//go:embed tmpl/corpus.go.tmpl
var tmplCorpus string

func captureCorpusValues(pkg *ast.Package, fname string, sourceCode *ast.File, fFuzz *ast.FuncLit) {
	tmpl, err := template.New(``).Parse(tmplCorpus)
	if err != nil {
		panic(err)
	}

	// Write corpus function to file
	tempFile, err := os.CreateTemp(`.`, `corpus-*_test.go`)
	if err != nil {
		panic(err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	err = tmpl.Execute(tempFile, pkg)
	if err != nil {
		panic(err)
	}

	// Call this function in f.Fuzz
	var paramNames []string
	for _, param := range fFuzz.Type.Params.List[1:] {
		paramNames = append(paramNames, param.Names[0].Name)
	}

	saveFunc, err := parser.ParseExpr(`_GoFuzzBuildSaveCorpus(` + strings.Join(paramNames, `,`) + `)`)
	if err != nil {
		panic(err)
	}
	fFuzz.Body.List = append(fFuzz.Body.List, nil)
	copy(fFuzz.Body.List[1:], fFuzz.Body.List)
	fFuzz.Body.List[0] = &ast.ExprStmt{X: saveFunc}

	f, err := os.OpenFile(fname, os.O_RDWR, 0o644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Rollback original file
	backup, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	defer func() {
		f.Truncate(0)
		f.WriteAt(backup, 0)
	}()

	// Overwrite original file
	f.Seek(0, 0)
	err = format.Node(f, token.NewFileSet(), sourceCode)
	if err != nil {
		panic(err)
	}

	// Run and save corpa
	err = os.MkdirAll(`corpus`, 0o755)
	if err != nil {
		panic(err)
	}
	b, err := exec.Command(`go`, `test`, `-run`, *funcName).CombinedOutput()
	if err != nil {
		fmt.Println(string(b))
		panic(err)
	}

	files, err := os.ReadDir(`corpus`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Added %d corpa\n", len(files))
}
