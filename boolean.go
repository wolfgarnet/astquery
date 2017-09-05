package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
)

// booleanQuery will determine if a part of the tree is solely booleans
type booleanQuery struct {
	depth      int
	expression ast.Expression
}

func (qo *booleanQuery) run(expression ast.Expression) error {
	qo.expression = expression
	ok := VerifyExpression(expression, qo.depth, newOnlyBooleanVerifier())
	if !ok {
		return fmt.Errorf("Expression does not only contain booleans")
	}

	return nil
}

func (qo *booleanQuery) get() ast.Expression {
	return qo.expression
}

// AcceptBoolean will only pass if the subtree is solely booleans.
func (q *Query) AcceptBoolean(depth int) *Query {
	q.operations = append(q.operations, &booleanQuery{
		depth: depth,
	})
	return q
}
