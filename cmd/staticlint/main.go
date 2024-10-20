// Package for static check
// Can be run by command go run ./cmd/staticlint/ ./...
package staticlint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"

	"honnef.co/go/tools/staticcheck"
)

func main() {
	mychecks := []*analysis.Analyzer{
		OsExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	multichecker.Main(
		mychecks...,
	)
}

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "os_exit",
	Doc:  "check for os.Exit",
	Run:  osExitAnalyzerRun,
}

func osExitAnalyzerRun(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if call, ok := node.(*ast.SelectorExpr); ok {
				if x, ok := call.X.(*ast.Ident); ok && x.Name == "os" && call.Sel.Name == "Exit" {
					pass.Reportf(call.Pos(), "os.Exit is forbidden")
				}
			}
			return true
		})
	}
	return nil, nil
}
