package astquery

import "github.com/robertkrimen/otto/ast"

// VerifyExpression will verify an expression given a depth and a verification function.
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

	case *ast.BooleanLiteral:
		return true

	case *ast.NumberLiteral:
		return true

	case *ast.StringLiteral:
		return true

	case *ast.UnaryExpression:
		return VerifyExpression(e.Operand, depth, verify)

	case *ast.VariableExpression:
		return VerifyExpression(e.Initializer, depth, verify)

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

		case *ast.BooleanLiteral:
			return false
		}

		return true
	}
}

func newOnlyBooleanVerifier() func(ast.Expression) bool {
	return func(e ast.Expression) bool {
		switch e.(type) {
		case *ast.StringLiteral:
			return false
		case *ast.NumberLiteral:
			return false
		}

		return true
	}
}

// Inspector interface for inspecting a given expression and is used for the inspect function.
type Inspector interface {
	Inspect(expression ast.Expression) Inspector
}

// CallInspector is an implementation of the Inspector interface for calls
type CallInspector struct {
	Call *ast.CallExpression
}

func (i *CallInspector) Inspect(expression ast.Expression) Inspector {
	call, isCall := expression.(*ast.CallExpression)
	if isCall {
		i.Call = call
	}
	return i
}

// Inspect will inspect a given expression, avoiding certain types of expressions.
func Inspect(node ast.Expression, inspector Inspector) Inspector {
	switch e := node.(type) {
	case *ast.VariableExpression:
		return Inspect(e.Initializer, inspector)

	default:
		return inspector.Inspect(node)
	}
}
