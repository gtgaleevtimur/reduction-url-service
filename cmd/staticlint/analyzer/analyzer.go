// Package analyzer ищет вызовы os.Exit в main packages и возвращает позицию согласно AST.
package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer - статический анализатор поиска вызова os.Exit.
var ExitAnalyzer = &analysis.Analyzer{
	Name: "myExitAnalyzer",
	Doc:  "Don't use os.Exit in main package",
	Run:  run,
}

// run - функция поиска статического анализатора os.Exit.
func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.SelectorExpr:
				if x.Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "os.Exit call in main package")

				}
			}
			return true
		})
	}
	return nil, nil
}
