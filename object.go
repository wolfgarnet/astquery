package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
)

type mustBeObjectLiteral struct {
	literal *ast.ObjectLiteral
}

func (qo *mustBeObjectLiteral) run(e ast.Expression) error {
	object, isObject := e.(*ast.ObjectLiteral)
	if !isObject {
		return fmt.Errorf("Not an object literal")
	}

	qo.literal = object
	return nil
}

func (qo *mustBeObjectLiteral) get() ast.Expression {
	return qo.literal
}

func (q *Query) MustBeObjectLiteral() *Query {
	q.operations = append(q.operations, &mustBeObjectLiteral{})
	return q
}
