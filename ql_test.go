package astquery

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/token"
	"testing"
)

func TestQL(t *testing.T) {
	left := &ast.NumberLiteral{
		Value:   1,
		Literal: "1",
	}
	/*
		right := &ast.NumberLiteral{
			Value:   1,
			Literal: "1",
		}
	*/
	right := &ast.StringLiteral{
		Value:   "1",
		Literal: "1",
	}
	binary := &ast.BinaryExpression{
		Operator: token.PLUS,
		Left:     left,
		Right:    right,
	}

	err := NewQuery().MustBeBinary().HasOperator(token.PLUS).OneSideOtherSide(NewQuery().AcceptNumbers(5), NewQuery().AcceptNumbers(5)).Run(binary)

	fmt.Printf("ERROR IS %v\n", err)
}

func TestQL2(t *testing.T) {
	left := &ast.NumberLiteral{
		Value:   1,
		Literal: "1",
	}

	r1 := &ast.NumberLiteral{
		Value:   2,
		Literal: "2",
	}
	r2 := &ast.NumberLiteral{
		Value:   2,
		Literal: "2",
	}
	right := &ast.BinaryExpression{
		Operator: token.PLUS,
		Left:     r1,
		Right:    r2,
	}
	binary := &ast.BinaryExpression{
		Operator: token.MULTIPLY,
		Left:     left,
		Right:    right,
	}



	err := NewQuery().MustBeBinary().HasOperator(token.MULTIPLY).OneSideOtherSide(
		NewQuery().AcceptNumbers(1),
		NewQuery().MustBeBinary().HasOperator(token.PLUS).AcceptNumbers(5)).Run(binary)

	fmt.Printf("ERROR IS %v\n", err)
}

func TestQL3(t *testing.T) {
	left := &ast.BooleanLiteral{
		Value:   true,
		Literal: "true",
	}

	right := &ast.BooleanLiteral{
		Value:   true,
		Literal: "true",
	}
	/*
	right := &ast.NumberLiteral{
		Value:   2,
		Literal: "2",
	}
	*/
	binary := &ast.BinaryExpression{
		Operator: token.LOGICAL_AND,
		Left:     left,
		Right:    right,
	}
	unary := &ast.UnaryExpression{
		Operator: token.NOT,
		Operand:binary,
	}

	err := NewQuery().MustBeUnary().HasOperator(token.NOT).Operands(NewQuery().MustBeBinary().HasOperator(token.LOGICAL_AND).Operands(NewQuery().AcceptBoolean(5))).Run(unary)

	fmt.Printf("ERROR IS %v\n", err)
}