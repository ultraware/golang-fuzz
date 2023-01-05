package main

import (
	"flag"
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

var (
	funcName  = flag.String("func", "", "name of the Fuzz function")
	genCorpus = flag.Bool("corpus", true, "generate corpus dir")
)

func main() {
	// Parse command line flags
	flag.Parse()
	if flag.NArg() == 0 || *funcName == "" {
		fmt.Println("Usage: go-fuzz-build -func FUNC_NAME PACKAGE_PATH")
		os.Exit(1)
	}
	packagePath := flag.Args()[0]

	// Find the Fuzz function in the package
	pkg, fname, sourceCode, fuzzFunc := findFuzzFunc(packagePath, *funcName)
	if fuzzFunc == nil {
		fmt.Printf("Fuzz function %s not found in package %s\n", *funcName, packagePath)
		os.Exit(1)
	}
	err := os.Chdir(packagePath)
	if err != nil {
		panic(err)
	}

	// Find the function passed to f.Fuzz
	fuzzCall := findFuzzCall(fuzzFunc.Body)
	if fuzzCall == nil {
		fmt.Printf("Fuzz function %s does not contain a call to f.Fuzz\n", *funcName)
		os.Exit(1)
	}

	fFuzz, ok := fuzzCall.Args[0].(*ast.FuncLit)
	if !ok {
		fmt.Printf("Expected function passed to f.Fuzz, got %T\n", fuzzCall.Args[0])
		os.Exit(1)
	}

	params := fFuzz.Type.Params.List
	fmt.Printf("Function passed to f.Fuzz has %d parameters:\n", len(params)-1)
	// Print the types of the parameters of the passed function
	for i := 1; i < len(params); i++ {
		fmt.Printf("- Parameter %d: %s\n", i, params[i].Type) // TODO: Check supported types
	}

	if !*genCorpus {
		return
	}
	captureCorpusValues(pkg, fname, sourceCode, fFuzz, params[1:])
}

func captureCorpusValues(pkg *ast.Package, fname string, sourceCode *ast.File, fFuzz *ast.FuncLit, params []*ast.Field) {
	tmpl, err := template.New(``).Parse(tmplCorpus)
	if err != nil {
		panic(err)
	}

	// Write corpus function to file
	tempFile, err := os.CreateTemp(`.`, `fuzz-*_test.go`)
	if err != nil {
		panic(err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	err = tmpl.Execute(tempFile, pkg)
	if err != nil {
		panic(err)
	}

	// Call fuction in f.Fuzz
	var paramNames []string
	for _, param := range params {
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

func findFuzzFunc(packagePath, funcName string) (*ast.Package, string, *ast.File, *ast.FuncDecl) {
	// Parse the Go package
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing package: %s\n", err)
		os.Exit(1)
	}

	// Find the Fuzz function in the package
	for _, pkg := range pkgs {
		for fname, file := range pkg.Files {
			for _, decl := range file.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == funcName {
					return pkg, fname, file, funcDecl
				}
			}
		}
	}
	return nil, ``, nil, nil
}

func findFuzzCall(node ast.Node) *ast.CallExpr {
	switch n := node.(type) {
	case *ast.ExprStmt:
		if call, ok := n.X.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.SelectorExpr); ok && ident.Sel.Name == "Fuzz" { // TODO: check if called by *testing.F
				return call
			}
		}
	case *ast.FuncLit:
		return findFuzzCall(n.Body)
	case *ast.BlockStmt:
		for _, stmt := range n.List {
			if result := findFuzzCall(stmt); result != nil {
				return result
			}
		}
	}
	return nil
}
