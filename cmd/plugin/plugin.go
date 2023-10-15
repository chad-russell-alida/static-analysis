// Plugin must be package main
package main

import (
	"github.com/golang-cz/static-analysis/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.NewWrapErrChecker(),
	}, nil
}
