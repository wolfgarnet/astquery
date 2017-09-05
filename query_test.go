package astquery

import (
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/token"
	"testing"
)

func TestQuery(t *testing.T) {
	tests := []struct {
		expression ast.Expression
		query      *Query
		mustFail   bool
	}{

		// Test 0
		{&ast.BinaryExpression{
			Operator: token.LOGICAL_OR,
			Left: &ast.BooleanLiteral{
				Value: true,
			},
			Right: &ast.BooleanLiteral{
				Value: true,
			},
		}, NewQuery().MustBeBinary().HasOperator(token.LOGICAL_OR).Operands(NewQuery().AcceptBoolean(5)),
			false,
		},

		// Test 1
		{&ast.BinaryExpression{
			Operator: token.LOGICAL_OR,
			Left: &ast.BooleanLiteral{
				Value: true,
			},
			Right: &ast.BooleanLiteral{
				Value: true,
			},
		}, NewQuery().MustBeUnary(),
			true,
		},

		// Test 2
		{&ast.BinaryExpression{
			Operator: token.LOGICAL_OR,
			Left: &ast.BooleanLiteral{
				Value: true,
			},
			Right: &ast.BooleanLiteral{
				Value: true,
			},
		}, NewQuery().MustBeBinary().Operands(NewQuery().AcceptNumbers(5)),
			true,
		},

		// Test 3
		{&ast.BinaryExpression{
			Operator: token.LOGICAL_OR,
			Left: &ast.BooleanLiteral{
				Value: true,
			},
			Right: &ast.NumberLiteral{
				Value: 5,
			},
		}, NewQuery().MustBeBinary().HasOperator(token.LOGICAL_OR).OneSideOtherSide(NewQuery().AcceptBoolean(5), NewQuery().AcceptNumbers(5)),
			false,
		},

		// Test 4
		{&ast.UnaryExpression{
			Operator: token.NOT,
			Operand: &ast.BooleanLiteral{
				Value: true,
			},
		}, NewQuery().MustBeUnary().HasOperator(token.NOT).Operands(NewQuery().AcceptBoolean(5)),
			false,
		},

		// Test 5
		{&ast.UnaryExpression{
			Operator: token.NOT,
			Operand: &ast.BooleanLiteral{
				Value: true,
			},
		}, NewQuery().MustBeBinary(),
			true,
		},

		// Test 6
		{&ast.UnaryExpression{
			Operator: token.INCREMENT,
			Operand: &ast.BooleanLiteral{
				Value: true,
			},
		}, NewQuery().MustBeUnary().HasOperator(token.NOT),
			true,
		},

		// Test 7
		{&ast.UnaryExpression{
			Operator: token.INCREMENT,
			Operand: &ast.NumberLiteral{
				Value: 5,
			},
		}, NewQuery().MustBeUnary().HasOperator(token.INCREMENT).Operands(NewQuery().AcceptNumbers(5)),
			false,
		},
	}

	for i, test := range tests {
		err := test.query.Run(test.expression)
		if err != nil && !test.mustFail {
			t.Errorf("Test %v failed, %v", i, err)
			continue
		}

		if err == nil && test.mustFail {
			t.Errorf("Test %v should have failed!", i)
			continue
		}
	}
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

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}
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
		Operand:  binary,
	}

	err := NewQuery().MustBeUnary().HasOperator(token.NOT).Operands(NewQuery().MustBeBinary().HasOperator(token.LOGICAL_AND).Operands(NewQuery().AcceptBoolean(5))).Run(unary)

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}
}

func TestQL4_fail(t *testing.T) {
	left := &ast.BooleanLiteral{
		Value:   true,
		Literal: "true",
	}

	right := &ast.NumberLiteral{
		Value:   2,
		Literal: "2",
	}

	binary := &ast.BinaryExpression{
		Operator: token.LOGICAL_AND,
		Left:     left,
		Right:    right,
	}
	unary := &ast.UnaryExpression{
		Operator: token.NOT,
		Operand:  binary,
	}

	err := NewQuery().MustBeUnary().HasOperator(token.NOT).Operands(NewQuery().MustBeBinary().HasOperator(token.LOGICAL_AND).Operands(NewQuery().AcceptBoolean(5))).Run(unary)

	if err == nil {
		t.Errorf("Test should have failed, %v", err)
	}
}

func TestQLCall(t *testing.T) {
	callee := &ast.Identifier{
		Name: "f",
	}
	call := &ast.CallExpression{
		Callee:       callee,
		ArgumentList: nil,
	}
	statement := &ast.ExpressionStatement{
		call,
	}

	err := NewQuery().MustBeCall().RunStatement(statement)

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}
}

func TestQLCall2(t *testing.T) {
	callee := &ast.Identifier{
		Name: "f",
	}
	call := &ast.CallExpression{
		Callee:       callee,
		ArgumentList: nil,
	}
	call2 := &ast.CallExpression{
		Callee:       call,
		ArgumentList: nil,
	}
	statement := &ast.ExpressionStatement{
		call2,
	}

	inspector := NewQuery().MustBeCallD(true).Collect()
	err := inspector.RunStatement(statement)

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}

	if call != inspector.Collected {
		t.Errorf("Collected expression not correct")
	}
}

func TestQLReturnStatement(t *testing.T) {
	callee := &ast.Identifier{
		Name: "f",
	}
	call := &ast.CallExpression{
		Callee:       callee,
		ArgumentList: nil,
	}

	statement := &ast.ReturnStatement{
		Argument: call,
	}

	err := NewQuery().MustBeCall().RunStatement(statement)

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}
}

func TestQLVarStatement(t *testing.T) {
	callee := &ast.Identifier{
		Name: "f",
	}
	call := &ast.CallExpression{
		Callee:       callee,
		ArgumentList: nil,
	}
	varExp := &ast.VariableExpression{
		Name:        "test",
		Initializer: call,
	}

	statement := &ast.VariableStatement{
		List: []ast.Expression{varExp},
	}

	err := NewQuery().MustBeCallD(false).RunStatement(statement)

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}
}

func TestQLAnonymousFunctionCall(t *testing.T) {

	function := &ast.FunctionLiteral{
		Body: nil,
	}

	call := &ast.CallExpression{
		Callee:       function,
		ArgumentList: nil,
	}

	statement := &ast.ExpressionStatement{
		Expression: call,
	}

	err := NewQuery().MustBeCall().RunStatement(statement)

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}
}

func TestQLAnonymousFunctionCall2(t *testing.T) {

	function := &ast.FunctionLiteral{
		Body: nil,
	}

	call := &ast.CallExpression{
		Callee:       function,
		ArgumentList: nil,
	}

	id := &ast.Identifier{
		Name: "i",
	}

	assign := &ast.AssignExpression{
		Left:  id,
		Right: call,
	}

	statement := &ast.ExpressionStatement{
		Expression: assign,
	}

	query := NewQuery().MustBeCallD(true).Collect()
	err := query.RunStatement(statement)

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}
}

func TestDotFunctionCall(t *testing.T) {

	object := &ast.Identifier{
		Name: "obj",
	}

	function := &ast.Identifier{
		Name: "f",
	}

	dot := &ast.DotExpression{
		Left:       object,
		Identifier: function,
	}

	call := &ast.CallExpression{
		Callee:       dot,
		ArgumentList: nil,
	}

	statement := &ast.ExpressionStatement{
		Expression: call,
	}

	q := NewQuery().MustBeCall().Collect()
	err := q.RunStatement(statement)

	if err != nil {
		t.Errorf("Test failed, %v", err)
	}
}

func TestEither(t *testing.T) {
	function := &ast.Identifier{
		Name: "f",
	}
	call := &ast.CallExpression{
		Callee:       function,
		ArgumentList: nil,
	}

	statement := &ast.ExpressionStatement{
		Expression: call,
	}

	err := NewQuery().Either(NewQuery().MustBeUnary(), NewQuery().MustBeBinary(), NewQuery().MustBeCall()).RunStatement(statement)
	if err != nil {
		t.Errorf("Failed: %v", err)
	}
}
