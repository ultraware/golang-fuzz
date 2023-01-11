package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func findFuzzFunc(packagePath, funcName string) (*ast.Package, string, *ast.FuncDecl) {
	// Parse the Go package
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packagePath, nil, 0)
	if err != nil {
		fmt.Printf("Error parsing package: %s\n", err)
		os.Exit(1)
	}

	// Find the Fuzz function in the package
	for _, pkg := range pkgs {
		for fname, file := range pkg.Files {
			if strings.HasSuffix(fname, `_test.go`) {
				continue
			}
			for _, decl := range file.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == funcName {
					return pkg, fname, funcDecl
				}
			}
		}
	}
	return nil, ``, nil
}
