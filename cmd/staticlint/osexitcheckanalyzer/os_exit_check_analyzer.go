// Package osexitchecker Запрещает использовать прямой вызов os.Exit в функции main пакета main
package osexitcheckanalyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// OsExitCheckAnalyzer Запрещает использовать прямой вызов os.Exit в функции main пакета main
var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "os_exit_check",
	Doc:  "check for os.Exit in main function in package main",
	Run:  run,
}

// Запуск анализатора
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if funDecl, ok := node.(*ast.FuncDecl); ok {
				if funDecl.Name.Name == "main" && file.Name.Name == "main" {
					ast.Inspect(funDecl, func(node ast.Node) bool {
						if exprStmt, ok := node.(*ast.ExprStmt); ok {
							if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
								if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
									if fun.Sel.Name == "Exit" && fun.X.(*ast.Ident).Name == "os" {
										pass.Reportf(fun.Pos(), "found os.Exit call")
									}
								}
							}
						}
						return true
					})
				}
			}

			return true
		})
	}
	return nil, nil
}
