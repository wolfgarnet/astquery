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
