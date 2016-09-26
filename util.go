package astquery

import (
	"github.com/robertkrimen/otto/ast"
)

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
	Done() bool
}

// CallInspector is an implementation of the Inspector interface for calls
type CallInspector struct {
	Call *ast.CallExpression
	New  *ast.NewExpression

	// First determines if the first found call is collected
	First bool
}

func (i *CallInspector) Inspect(expression ast.Expression) Inspector {
	call, isCall := expression.(*ast.CallExpression)
	if isCall && ((i.Call == nil && !i.First) || (i.First)) {
		i.Call = call
	}
	newe, isNew := expression.(*ast.NewExpression)
	if isNew && ((i.New == nil && !i.First) || (i.First)) {
		i.New = newe
	}
	return i
}

func (i *CallInspector) Done() bool {
	if i.First {
		return false
	} else {
		return i.Call != nil || i.New != nil
	}
}

// Inspect will inspect a given expression, avoiding certain types of expressions.
func Inspect(node ast.Expression, inspector Inspector) Inspector {

	if inspector.Inspect(node).Done() {
		return inspector
	}

	switch e := node.(type) {
	case *ast.ArrayLiteral:
		for _, v := range e.Value {
			if Inspect(v, inspector).Done() {
				return inspector
			}
		}
	case *ast.AssignExpression:
		if Inspect(e.Left, inspector).Done() {
			return inspector
		}
		if Inspect(e.Right, inspector).Done() {
			return inspector
		}
	case *ast.BinaryExpression:
		if Inspect(e.Left, inspector).Done() {
			return inspector
		}
		if Inspect(e.Right, inspector).Done() {
			return inspector
		}
	case *ast.BracketExpression:
		if Inspect(e.Left, inspector).Done() {
			return inspector
		}
		if Inspect(e.Member, inspector).Done() {
			return inspector
		}
	case *ast.CallExpression:
		if Inspect(e.Callee, inspector).Done() {
			return inspector
		}
		for _, l := range e.ArgumentList {
			if Inspect(l, inspector).Done() {
				return inspector
			}
		}
	case *ast.ConditionalExpression:
		if Inspect(e.Test, inspector).Done() {
			return inspector
		}
		if Inspect(e.Consequent, inspector).Done() {
			return inspector
		}
		if Inspect(e.Alternate, inspector).Done() {
			return inspector
		}
	case *ast.DotExpression:
		if Inspect(e.Left, inspector).Done() {
			return inspector
		}
		if Inspect(e.Identifier, inspector).Done() {
			return inspector
		}
	case *ast.FunctionLiteral:
		if Inspect(e.Name, inspector).Done() {
			return inspector
		}
		for _, l := range e.ParameterList.List {
			if Inspect(l, inspector).Done() {
				return inspector
			}
		}
	case *ast.NewExpression:
		if Inspect(e.Callee, inspector).Done() {
			return inspector
		}
		for _, l := range e.ArgumentList {
			if Inspect(l, inspector).Done() {
				return inspector
			}
		}
	case *ast.ObjectLiteral:
		for _, l := range e.Value {
			if Inspect(l.Value, inspector).Done() {
				return inspector
			}
		}
	case *ast.SequenceExpression:
		for _, l := range e.Sequence {
			if Inspect(l, inspector).Done() {
				return inspector
			}
		}
	case *ast.UnaryExpression:
		if Inspect(e.Operand, inspector).Done() {
			return inspector
		}
	case *ast.VariableExpression:
		if Inspect(e.Initializer, inspector).Done() {
			return inspector
		}
	}

	return inspector
}
