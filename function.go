package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
)

type functionLiteralQuery struct {
	function ast.Expression
}

func (qo *functionLiteralQuery) run(e ast.Expression) error {
	var isFLiteral bool
	qo.function, isFLiteral = e.(*ast.FunctionLiteral)
	if !isFLiteral {
		return fmt.Errorf("Expression is not a function literal")
	}

	return nil
}

func (qo *functionLiteralQuery) get() ast.Expression {
	return qo.function
}

func (q *Query) MustBeFunctionLiteral() *Query {
	q.operations = append(q.operations, &functionLiteralQuery{})
	return q
}
