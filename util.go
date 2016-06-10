package astquery

import "github.com/robertkrimen/otto/ast"

func VerifyExpression(exp ast.Expression, depth int, verify func(ast.Expression) bool) bool {

	if depth == 0 {
		return false
	}

	depth--

	if !verify(exp) {
		return false
	}

	switch e := exp.(type) {
	case *ast.BinaryExpression:
		if !VerifyExpression(e.Left, depth, verify) {
			return false
		}
		if !VerifyExpression(e.Right, depth, verify) {
			return false
		}

	case *ast.NumberLiteral:
		return true

	case *ast.StringLiteral:
		return true

	default:
		return false
	}

	return true
}

func newOnlyNumberVerifier() func(ast.Expression) bool {
	return func(e ast.Expression) bool {
		switch e.(type) {
		case *ast.StringLiteral:
			return false
		}

		return true
	}
}
