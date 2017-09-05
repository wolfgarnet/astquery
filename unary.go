package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
)

// unaryQuery requires the current expression to be unary
type unaryQuery struct {
	unary *ast.UnaryExpression
}

func (qo *unaryQuery) run(e ast.Expression) error {
	var isUnary bool
	qo.unary, isUnary = e.(*ast.UnaryExpression)
	if !isUnary {
		return fmt.Errorf("Expression is not unary, was %T", e)
	}

	return nil
}

func (qo *unaryQuery) get() ast.Expression {
	return qo.unary
}

// MustBeUnary restricts the expression to be unary
func (q *Query) MustBeUnary() *Query {
	q.operations = append(q.operations, &unaryQuery{})
	return q
}

