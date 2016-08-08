package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/token"
)

// Query defines the basic structure for a query
type Query struct {
	operations []QLOperation
	Collected  ast.Expression
}

// NewQuery returns a new query
func NewQuery() *Query {
	return &Query{}
}

// Run runs the query given the ast expression
func (ql *Query) Run(expression ast.Expression) error {
	for _, q := range ql.operations {
		err := q.run(expression)
		if err != nil {
			return err
		}

		expression = q.get()
		if ql.Collected == nil {
			ql.Collected = expression
		}
	}

	return nil
}

func (ql *Query) RunStatement(statement ast.Statement) error {
	switch s := statement.(type) {
	case *ast.ExpressionStatement:
		return ql.Run(s.Expression)
	case *ast.ReturnStatement:
		return ql.Run(s.Argument)
	case *ast.VariableStatement:
		for _, e := range s.List {
			err := ql.Run(e)
			if err != nil {
				return err
			}
		}

		return nil
	default:
		return fmt.Errorf("Unsupported statement: %T", statement)
	}
}

func (ql *Query) Collect() *Query {
	ql.Collected = nil
	return ql
}

// QLOperation specifies a query operation
type QLOperation interface {
	run(ast.Expression) error
	get() ast.Expression
}

// binaryQuery requires the current expression to be binary
type binaryQuery struct {
	binary *ast.BinaryExpression
}

func (qo *binaryQuery) run(e ast.Expression) error {
	var isBinary bool
	qo.binary, isBinary = e.(*ast.BinaryExpression)
	if !isBinary {
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

// callQuery
type callQuery struct {
	depth int
	call  *ast.CallExpression
}

func (qo *callQuery) run(e ast.Expression) error {
	if qo.depth > 0 {
		inspector := &CallInspector{}
		Inspect(e, inspector)
		ok := inspector.Call != nil
		if !ok {
			return fmt.Errorf("Expression does not contain a call")
		}
		qo.call = inspector.Call
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
	return qo.call
}

// MustBeBinary restricts the expression to be binary
func (q *Query) MustBeCall() *Query {
	q.operations = append(q.operations, &callQuery{})
	return q
}

func (q *Query) MustBeCallD() *Query {
	q.operations = append(q.operations, &callQuery{1, nil})
	return q
}

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

// operatorQuery filters expressions based on operators
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

// HasOperator filters expressions given the set of operators.
func (q *Query) HasOperator(operators ...token.Token) *Query {
	q.operations = append(q.operations, &operatorQuery{
		operators: operators,
	})
	return q
}

// eitherSideQuery will try to run the two provided queries on the binary expression in both order.
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

// OneSideOtherSide will run the queries on both operands in a binary expression in both order.
// If the first order doesn't work the other is tried.
func (q *Query) OneSideOtherSide(one *Query, other *Query) *Query {
	q.operations = append(q.operations, &eitherSideQuery{
		one:   one,
		other: other,
	})
	return q
}

// numberQuery will determine if a part of the tree is solely numbers
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

// AcceptNumbers will only pass if the subtree is solely numbers.
func (q *Query) AcceptNumbers(depth int) *Query {
	q.operations = append(q.operations, &numberQuery{
		depth: depth,
	})
	return q
}

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

// operandsQuery will run a query on all possible operands
type operandsQuery struct {
	expression ast.Expression
	query      *Query
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

// Operands will run the query on all possible operands.
// Unary - one operand
// Binary - two operands
func (q *Query) Operands(query *Query) *Query {
	q.operations = append(q.operations, &operandsQuery{
		query: query,
	})
	return q
}
