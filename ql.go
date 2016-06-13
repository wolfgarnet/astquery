package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/token"
)

type Query struct {
	operations []QLOperation
}

func NewQuery() *Query {
	return &Query{}
}

func (ql *Query) Run(expression ast.Expression) error {
	for _, q := range ql.operations {
		err := q.run(expression)
		if err != nil {
			return err
		}

		expression = q.get()
	}

	return nil
}

type QLOperation interface {
	run(ast.Expression) error
	get() ast.Expression
}

type binaryQuery struct {
	binary *ast.BinaryExpression
}

func (qo *binaryQuery) run(e ast.Expression) error {
	var isBinary bool
	qo.binary, isBinary = e.(*ast.BinaryExpression)
	if !isBinary {
		return fmt.Errorf("Expression is not binary")
	}

	return nil
}

func (qo *binaryQuery) get() ast.Expression {
	return qo.binary
}

func (q *Query) MustBeBinary() *Query {
	q.operations = append(q.operations, &binaryQuery{})
	return q
}


type unaryQuery struct {
	unary *ast.UnaryExpression
}

func (qo *unaryQuery) run(e ast.Expression) error {
	var isUnary bool
	qo.unary, isUnary = e.(*ast.UnaryExpression)
	if !isUnary {
		return fmt.Errorf("Expression is not unary")
	}

	return nil
}

func (qo *unaryQuery) get() ast.Expression {
	return qo.unary
}

func (q *Query) MustBeUnary() *Query {
	q.operations = append(q.operations, &unaryQuery{})
	return q
}


type operatorQuery struct {
	operators  []token.Token
	expression ast.Expression
}

func (qo *operatorQuery) run(e ast.Expression) error {
	qo.expression = e
	var operator token.Token
	switch t := e.(type) {
	case *ast.BinaryExpression:
		operator = t.Operator

	case *ast.UnaryExpression:
		operator = t.Operator

	default:
		return fmt.Errorf("Expression not compatible with operators")
	}

	found := false
	for _, op := range qo.operators {
		if op == operator {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Invalid operator for expression, %v", operator)
	}

	return nil
}

func (qo *operatorQuery) get() ast.Expression {
	return qo.expression
}

func (q *Query) HasOperator(operators ...token.Token) *Query {
	q.operations = append(q.operations, &operatorQuery{
		operators: operators,
	})
	return q
}

type eitherSideQuery struct {
	expression ast.Expression
	one, other *Query
}

func (qo *eitherSideQuery) run(e ast.Expression) error {

	// Must be binary
	binary, isBinary := e.(*ast.BinaryExpression)
	if !isBinary {
		return fmt.Errorf("Expression is not binary")
	}

	// First
	err1 := qo.one.Run(binary.Left)
	err2 := qo.other.Run(binary.Right)
	if err1 != nil || err2 != nil {
		err1 := qo.one.Run(binary.Right)
		err2 := qo.other.Run(binary.Left)

		if err1 != nil || err2 != nil {
			return fmt.Errorf("Expression is not compatible, left: %v, right: %v", err1, err2)
		}
	}

	qo.expression = e

	return nil
}

func (qo *eitherSideQuery) get() ast.Expression {
	return qo.expression
}

func (q *Query) OneSideOtherSide(one *Query, other *Query) *Query {
	q.operations = append(q.operations, &eitherSideQuery{
		one:   one,
		other: other,
	})
	return q
}

type numberQuery struct {
	depth      int
	expression ast.Expression
}

func (qo *numberQuery) run(expression ast.Expression) error {
	qo.expression = expression
	ok := VerifyExpression(expression, qo.depth, newOnlyNumberVerifier())
	if !ok {
		return fmt.Errorf("Expression does not only contain numbers")
	}

	return nil
}

func (qo *numberQuery) get() ast.Expression {
	return qo.expression
}

func (q *Query) AcceptNumbers(depth int) *Query {
	q.operations = append(q.operations, &numberQuery{
		depth: depth,
	})
	return q
}

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

func (q *Query) AcceptBoolean(depth int) *Query {
	q.operations = append(q.operations, &booleanQuery{
		depth: depth,
	})
	return q
}

type operandsQuery struct {
	expression ast.Expression
	query *Query
}

func (qo *operandsQuery) run(expression ast.Expression) error {
	switch t := expression.(type) {
	case *ast.BinaryExpression:
		err1 := qo.query.Run(t.Left)
		err2 := qo.query.Run(t.Right)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("Binary operands where not compatible, left: %v, right: %v", err1, err2)
		}
	case *ast.UnaryExpression:
		err := qo.query.Run(t.Operand)
		return err
	default:
		return fmt.Errorf("Expression does not have one operand")
	}

	return nil
}

func (qo *operandsQuery) get() ast.Expression {
	return qo.expression
}

func (q *Query) Operands(query *Query) *Query {
	q.operations = append(q.operations, &operandsQuery{
		query:query,
	})
	return q
}
