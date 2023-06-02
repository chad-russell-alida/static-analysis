package analyzer

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

const (
	name = "wraperrchecker"
	doc  = "Checks wrapping of errors"
)

func NewWrapErrChecker() *analysis.Analyzer {
	n := newWrapErrChecker()

	a := &analysis.Analyzer{
		Name:     name,
		Doc:      doc,
		Run:      n.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	return a
}

type wrapErrChecker struct{}

func newWrapErrChecker() *wrapErrChecker {
	return &wrapErrChecker{}
}

func (*wrapErrChecker) run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.Nodes([]ast.Node{(*ast.IfStmt)(nil)}, func(node ast.Node, push bool) (proceed bool) {
		switch ifStmt := node.(type) {
		case *ast.IfStmt:
			binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
			if !ok {
				return true
			}

			if binExpr.Op != token.NEQ {
				return true
			}

			errIdent, ok := binExpr.X.(*ast.Ident)
			if !ok || errIdent.Name != "err" {
				return true
			}

			ast.Inspect(ifStmt.Body, func(n ast.Node) bool {
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				funSelExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				if funSelExpr.Sel.Name != "Errorf" {
					return true
				}

				selector, ok := funSelExpr.X.(*ast.Ident)
				if !ok || selector.Name != "fmt" {
					return true
				}

				argExprs := callExpr.Args
				if len(argExprs) < 1 {
					pass.Reportf(ifStmt.Pos(), "invalid number of arguments")
				} else {
					formatExpr := argExprs[0]
					formatString := extractFormatString(formatExpr)
					if formatString == "" {
						pass.Reportf(ifStmt.Pos(), "format string in fmt.Errorf should contain wrapped error")
					} else if !hasVerb(formatString, "%w") {
						pass.Reportf(ifStmt.Pos(), "format string in fmt.Errorf should contain wrapped error")
					}
				}

				return true
			})
		}

		return true
	})

	return nil, nil
}

// extractFormatString extracts the format string from the argument expression of fmt.Errorf.
func extractFormatString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.BasicLit:
		if t.Kind == token.STRING {
			return strings.Trim(t.Value, "\"`")
		}
	case *ast.BinaryExpr:
		return extractFormatString(t.X) + extractFormatString(t.Y)
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return extractFormatString(t.Sel)
	case *ast.CallExpr:
		args := make([]string, len(t.Args))
		for i, arg := range t.Args {
			args[i] = extractFormatString(arg)
		}
		return strings.Join(args, "")
	case *ast.StarExpr:
		return extractFormatString(t.X)
	}

	return ""
}

// hasVerb checks if the format string contains the specified verb.
func hasVerb(formatString, verb string) bool {
	return strings.Contains(formatString, verb)
}
