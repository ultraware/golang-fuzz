package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

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