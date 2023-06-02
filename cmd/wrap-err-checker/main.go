package main

import (
	"github.com/golang-cz/static-analysis/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.NewWrapErrChecker())
}
