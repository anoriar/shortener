package main

import (
	"github.com/anoriar/shortener/cmd/staticlint/osexitcheckanalyzer"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

// go -o build mycheck
// ./mycheck ./multicheckertestdata
func main() {
	multichecker.Main(
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		shift.Analyzer,
		osexitcheckanalyzer.OsExitCheckAnalyzer,
	)
}
