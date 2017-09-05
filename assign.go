package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
)

// assignQuery requires the current expression to be assign
type assignQuery struct {
	assign ast.Expression
}

func (qo *assignQuery) run(e ast.Expression) error {
	var isAssign bool
	qo.assign, isAssign = e.(*ast.AssignExpression)
	if !isAssign {
		return fmt.Errorf("Expression is not assign, was %T", e)
	}

	return nil
}

func (qo *assignQuery) get() ast.Expression {
	return qo.assign
}

// MustBeAssign restricts the expression to be binary
func (q *Query) MustBeAssign() *Query {
	q.operations = append(q.operations, &assignQuery{})
	return q
}


type assignOrVarQuery struct {
	assignOrVar ast.Expression
}

func (qo *assignOrVarQuery) run(e ast.Expression) error {
	var isAssign, isVar bool
	qo.assignOrVar, isAssign = e.(*ast.AssignExpression)
	if !isAssign {
		qo.assignOrVar, isVar = e.(*ast.VariableExpression)
		if !isVar {
			return fmt.Errorf("Expression is not a variable or assign expression, was %T", e)
		}
	}

	return nil
}

func (qo *assignOrVarQuery) get() ast.Expression {
	return qo.assignOrVar
}

// MustBeAssign restricts the expression to be binary
func (q *Query) MustBeAssignOrVar() *Query {
	q.operations = append(q.operations, &assignOrVarQuery{})
	return q
}
