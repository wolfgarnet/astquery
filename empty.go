package astquery

import "github.com/robertkrimen/otto/ast"

// Empty
type emptyQuery struct {
	expression ast.Expression
}

func (qo *emptyQuery) run(e ast.Expression) error {
	qo.expression = e
	return nil
}

func (qo *emptyQuery) get() ast.Expression {
	return qo.expression
}

// Empty does nothing, but collects the expression(if needed)
func (q *Query) Empty() *Query {
	q.operations = append(q.operations, &emptyQuery{})
	return q
}
