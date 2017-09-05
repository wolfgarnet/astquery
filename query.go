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

// operatorQuery filters expressions based on operators
type operatorQuery struct {
	operators  []token.Token
	expression ast.Expression
}

func (qo *operatorQuery) run(e ast.Expression) error {
	qo.expression = e
	var operator token.Token
	switch t := e.(type) {
	case *ast.AssignExpression:
		operator = t.Operator

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

func (q *Query) RightSide(query *Query) *Query {
	q.operations = append(q.operations, &rightSideQuery{
		query: query,
	})
	return q
}

type rightSideQuery struct {
	expression ast.Expression
	query      *Query
}

func (qo *rightSideQuery) run(e ast.Expression) error {
	qo.expression = e
	switch n := e.(type) {
	case *ast.AssignExpression:
		return qo.query.Run(n.Right)
	case *ast.BinaryExpression:
		return qo.query.Run(n.Right)
	case *ast.VariableExpression:
		return qo.query.Run(n.Initializer)
	default:
		return fmt.Errorf("Expression is not compatible with right side queries.")
	}
}

func (qo *rightSideQuery) get() ast.Expression {
	return qo.expression
}

type either struct {
	expression ast.Expression
	queries    []*Query
}

func (qo *either) run(e ast.Expression) error {
	errors := make([]error, len(qo.queries))
	failed := true
	for i, q := range qo.queries {
		err := q.Run(e)
		if err == nil {
			failed = false
			break
		}

		errors[i] = err
	}

	if failed {
		return fmt.Errorf("Failed either: %v", errors)
	}

	qo.expression = e
	return nil
}

func (qo *either) get() ast.Expression {
	return qo.expression
}

func (q *Query) Either(queries ...*Query) *Query {
	q.operations = append(q.operations, &either{
		queries: queries,
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
