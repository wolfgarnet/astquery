package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
)

// binaryQuery requires the current expression to be binary
type binaryQuery struct {
	binary ast.Expression
}

func (qo *binaryQuery) run(e ast.Expression) error {
	switch e.(type) {
	case *ast.AssignExpression:
		qo.binary = e
	case *ast.BinaryExpression:
		qo.binary = e
	default:
		return fmt.Errorf("Expression is not binary, was %T", e)
	}

	return nil
}

func (qo *binaryQuery) get() ast.Expression {
	return qo.binary
}

// MustBeBinary restricts the expression to be binary
func (q *Query) MustBeBinary() *Query {
	q.operations = append(q.operations, &binaryQuery{})
	return q
}

