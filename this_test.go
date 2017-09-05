package astquery

import (
	"testing"
	"github.com/robertkrimen/otto/ast"
)

func TestQuery_ContainsThis_no_this(t *testing.T) {
	identifier := &ast.Identifier{
		Name:"test",
	}

	statement := &ast.VariableStatement{
		List: []ast.Expression{identifier},
	}

	err := NewQuery().ContainsThis().RunStatement(statement)

	if err == nil {
		t.Errorf("Test failed, %v", err)
	}
}

func TestQuery_ContainsThis(t *testing.T) {
	identifier := &ast.Identifier{
		Name:"test",
	}

	assign := &ast.AssignExpression{
		Left:identifier,
		Right:&ast.ThisExpression{},
	}

	statement := &ast.VariableStatement{
		List: []ast.Expression{assign},
	}

	err := NewQuery().ContainsThis().RunStatement(statement)

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}
}

func TestQuery_ContainsThis2(t *testing.T) {
	identifier := &ast.Identifier{
		Name:"test",
	}

	binary := &ast.BinaryExpression{
		Left:&ast.NumberLiteral{
			Value:1,
			Literal:"1",
		},
		Right:&ast.BinaryExpression{
			Left:&ast.NumberLiteral{
				Value:2,
				Literal:"2",
			},
			Right:&ast.DotExpression{
				Left:&ast.ThisExpression{},
				Identifier:&ast.Identifier{
					Name:"test2",
				},
			},
		},
	}

	assign := &ast.AssignExpression{
		Left:identifier,
		Right:binary,
	}

	statement := &ast.VariableStatement{
		List: []ast.Expression{assign},
	}

	err := NewQuery().ContainsThis().RunStatement(statement)

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}
}
