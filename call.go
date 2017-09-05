package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
)

// callQuery
type callQuery struct {
	depth int
	first bool
	call  *ast.CallExpression
	newe  *ast.NewExpression
}

func (qo *callQuery) run(e ast.Expression) error {
	if qo.depth > 0 {
		inspector := &CallInspector{}
		inspector.First = qo.first
		Inspect(e, inspector)
		ok := inspector.Call != nil || inspector.New != nil
		if !ok {
			return fmt.Errorf("Expression does not contain a call")
		}
		qo.call = inspector.Call
		qo.newe = inspector.New
	} else {
		var isCall bool
		qo.call, isCall = e.(*ast.CallExpression)
		if !isCall {
			return fmt.Errorf("Expression is not call, was %T", e)
		}
	}

	return nil
}

func (qo *callQuery) get() ast.Expression {
	if qo.call != nil {
		return qo.call
	}

	return qo.newe
}

func (q *Query) MustBeAnonymousCall() *Query {
	q.operations = append(q.operations, &callQuery{})
	return q
}

// MustBeCall restricts the expression to be a call
func (q *Query) MustBeCall() *Query {
	q.operations = append(q.operations, &callQuery{})
	return q
}

func (q *Query) MustBeCallD(first bool) *Query {
	q.operations = append(q.operations, &callQuery{1, first, nil, nil})
	return q
}

// Callee name
type calleeName struct {
	identifier *ast.Identifier
}

func (qo *calleeName) run(e ast.Expression) error {
	call, isCall := e.(*ast.CallExpression)
	if !isCall {
		return fmt.Errorf("Expression is not a call")
	}

	switch t := call.Callee.(type) {
	case *ast.Identifier:
		qo.identifier = t
	case *ast.DotExpression:
		qo.identifier = t.Identifier

	default:
		return fmt.Errorf("Call expression does not contain an identifier")
	}

	return nil
}

func (qo *calleeName) get() ast.Expression {
	return qo.identifier
}

func (q *Query) CallMustHaveIdentifier() *Query {
	q.operations = append(q.operations, &calleeName{})
	return q
}
