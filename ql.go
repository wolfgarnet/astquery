package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/token"
)

type QL struct {
	operations []QLOperation
}

func NewQuery() *QL {
	return &QL{}
}

func (ql *QL) Run(expression ast.Expression) error {
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

type QLBinary struct {
	binary *ast.BinaryExpression
}

func (o *QLBinary) run(e ast.Expression) error {
	var isBinary bool
	o.binary, isBinary = e.(*ast.BinaryExpression)
	if !isBinary {
		return fmt.Errorf("Expression is not binary")
	}

	return nil
}

func (o *QLBinary) get() ast.Expression {
	return o.binary
}

func (ql *QL) MustBeBinary() *QL {
	ql.operations = append(ql.operations, &QLBinary{})
	return ql
}


type QLUnary struct {
	unary *ast.UnaryExpression
}

func (o *QLUnary) run(e ast.Expression) error {
	var isUnary bool
	o.unary, isUnary = e.(*ast.UnaryExpression)
	if !isUnary {
		return fmt.Errorf("Expression is not unary")
	}

	return nil
}

func (o *QLUnary) get() ast.Expression {
	return o.unary
}

func (ql *QL) MustBeUnary() *QL {
	ql.operations = append(ql.operations, &QLUnary{})
	return ql
}


type QLOperator struct {
	operators  []token.Token
	expression ast.Expression
}

func (o *QLOperator) run(e ast.Expression) error {
	o.expression = e
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
	for _, op := range o.operators {
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

func (o *QLOperator) get() ast.Expression {
	return o.expression
}

func (ql *QL) HasOperator(operators ...token.Token) *QL {
	ql.operations = append(ql.operations, &QLOperator{
		operators: operators,
	})
	return ql
}

type QLSideOther struct {
	expression ast.Expression
	one, other *QL
}

func (o *QLSideOther) run(e ast.Expression) error {

	// Must be binary
	binary, isBinary := e.(*ast.BinaryExpression)
	if !isBinary {
		return fmt.Errorf("Expression is not binary")
	}

	// First
	err1 := o.one.Run(binary.Left)
	err2 := o.other.Run(binary.Right)
	if err1 != nil || err2 != nil {
		err1 := o.one.Run(binary.Right)
		err2 := o.other.Run(binary.Left)

		if err1 != nil || err2 != nil {
			return fmt.Errorf("Expression is not compatible, left: %v, right: %v", err1, err2)
		}
	}

	o.expression = e

	return nil
}

func (o *QLSideOther) get() ast.Expression {
	return o.expression
}

func (ql *QL) OneSideOtherSide(one *QL, other *QL) *QL {
	ql.operations = append(ql.operations, &QLSideOther{
		one:   one,
		other: other,
	})
	return ql
}

type QLNumberLeaf struct {
	depth      int
	expression ast.Expression
}

func (o *QLNumberLeaf) run(expression ast.Expression) error {
	o.expression = expression
	ok := VerifyExpression(expression, o.depth, newOnlyNumberVerifier())
	if !ok {
		return fmt.Errorf("Expression does not only contain numbers")
	}

	return nil
}

func (o *QLNumberLeaf) get() ast.Expression {
	return o.expression
}

func (ql *QL) AcceptNumbers(depth int) *QL {
	ql.operations = append(ql.operations, &QLNumberLeaf{
		depth: depth,
	})
	return ql
}

type QLOperand struct {
	expression ast.Expression
}

func (o *QLOperand) run(expression ast.Expression) error {
	switch t := expression.(type) {
	case *ast.UnaryExpression:
		o.expression = t.Operand
	default:
		return fmt.Errorf("Expression does not have one operand")
	}

	return nil
}

func (o *QLOperand) get() ast.Expression {
	return o.expression
}

func (ql *QL) Operand() *QL {
	ql.operations = append(ql.operations, &QLOperand{})
	return ql
}
