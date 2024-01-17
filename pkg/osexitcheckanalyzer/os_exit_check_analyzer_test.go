package osexitcheckanalyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitCheckerAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OsExitCheckAnalyzer, "./...")
}
