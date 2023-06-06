//go:build tools

// This must be package main
package main

import (
	"github.com/golang-cz/static-analysis/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// This must be defined and named 'AnalyzerPlugin'
var AnalyzerPlugin analyzerPlugin

type analyzerPlugin struct{}

// This must be implemented
func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{analyzer.NewWrapErrChecker()}
}
