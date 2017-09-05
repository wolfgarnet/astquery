package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
)

// booleanQuery will determine if a part of the tree is solely booleans
type thisQuery struct {
	expression ast.Expression
}

func (qo *thisQuery) run(expression ast.Expression) error {
	inspector := &ThisInspector{}
	Inspect(expression, inspector)

	if inspector.Found == 0 {
		return fmt.Errorf("Expression does not contain this")
	}

	return nil
}

func (qo *thisQuery) get() ast.Expression {
	return qo.expression
}

// AcceptBoolean will only pass if the subtree is solely booleans.
func (q *Query) ContainsThis() *Query {
	q.operations = append(q.operations, &thisQuery{
	})
	return q
}
