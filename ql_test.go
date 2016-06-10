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

	err := NewQuery().MustBeBinaryExpression().HasOperator(token.PLUS).OneSideOtherSide(NewQuery().AcceptNumbers(5), NewQuery().AcceptNumbers(5)).Run(binary)

	fmt.Printf("ERROR IS %v\n", err)
}
