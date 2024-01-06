package osexitcheckanalyzer

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestOsExitCheckerAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OsExitCheckAnalyzer, "./...")
}
