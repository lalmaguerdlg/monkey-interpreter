package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		testLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedValue)
	}
}

func TestAssignmentStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"x = 5;", "x", 5},
		{"y = true;", "y", true},
		{"foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		stmt := program.Statements[0]
		testAssignmentStatement(t, stmt, tt.expectedIdentifier, tt.expectedValue)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{"return x;", "x"},
		{"return 5;", 5},
		{"return true;", true},
		// TODO: add tests for generic ast.Expression
		// Create a helper function testExpression that compares 2 ast.Expression and they must be equal expressions
		//
		// {"return a + b;", &ast.InfixExpression{
		// 	Token:    token.Token{Type: token.PLUS, Literal: "+"},
		// 	Operator: "+",
		// 	Left: &ast.Identifier{
		// 		Token: token.Token{Type: token.IDENT, Literal: "a"},
		// 		Value: "a",
		// 	},
		// 	Right: &ast.Identifier{
		// 		Token: token.Token{Type: token.IDENT, Literal: "b"},
		// 		Value: "b",
		// 	},
		// }},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("program.Statements[0] not *ast.ReturnStatement. got=%T", returnStmt)
			continue
		}
		if returnStmt.TokenType() != token.RETURN {
			t.Errorf("returnStmt.TokenType not %s. got=%s", token.RETURN, string(returnStmt.TokenType()))
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return'. got=%s", returnStmt.TokenLiteral())
		}

		testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue)
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string, value any) bool {
	if string(stmt.TokenType()) != token.LET {
		t.Errorf("stmt.TokenType not '%q'. got=%q", token.LET, stmt.TokenLiteral())
		return false
	}
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not *ast.LetStatement. got=%T", stmt)
		return false
	}

	if string(letStmt.Name.TokenType()) != token.IDENT {
		t.Errorf("letStmt.Name.TokeType not '%q'. got=%q", token.IDENT, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%q", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return testLiteralExpression(t, letStmt.Value, value)
}

func testAssignmentStatement(t *testing.T, stmt ast.Statement, name string, value any) bool {
	assignment, ok := stmt.(*ast.AssignmentStatement)
	if !ok {
		t.Errorf("stmt not *ast.AssignmentStatement. got=%T", stmt)
		return false
	}

	testIdentifier(t, assignment.Name, name)

	return testLiteralExpression(t, assignment.Value, value)
}

func TestIdentifierExpressions(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if !testIdentifier(t, stmt.Expression, "foobar") {
		return
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if !testIntegerLiteral(t, stmt.Expression, 5) {
		return
	}
}

func TestBooleanExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. expected=%d got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testBooleanLiteral(t, stmt.Expression, tt.expected) {
			return
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world";`, "hello world"},
		{`"hello \"world\"";`, `hello "world"`},
		{`"hello\nworld";`, "hello\nworld"},
		{`"hello\tworld";`, "hello\tworld"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testStringLiteral(t, stmt.Expression, tt.expected) {
			return
		}
	}
}

func TestIfExpressions(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("exp.Consequence.Statements does not contain %d statements. got=%d\n", 1, len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Fatalf("exp.Alternative was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpressions(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("exp.Consequence.Statements does not contain %d statements. got=%d\n", 1, len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiterals(t *testing.T) {
	input := `fn(x, y, z) { x + y + z; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 3 {
		t.Fatalf("function definition does not contain %d parameters. got=%d\n", 3, len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	testLiteralExpression(t, function.Parameters[2], "z")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements does not contain %d statements. got=%d\n", 1, len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	expectedString := "((x + y) + z)"
	if bodyStmt.String() != expectedString {
		t.Fatalf("function body stmt is did not match expected string representation. got=%s expected=%s", bodyStmt.String(), expectedString)
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y) {};", []string{"x", "y"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Errorf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
		}

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("function definition does not contain %d parameters. got=%d\n", len(tt.expectedParams), len(function.Parameters))
		}

		for i, param := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], param)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	tests := []struct {
		input             string
		expectedFunction  string
		expectedArguments any
	}{
		{"run();", "run", []string{}},
		{"add(1, 2)", "add", []any{1, 2}},
		{
			"add(1, 2 + 2)",
			"add",
			[]ast.Expression{
				&ast.IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "1"},
					Value: 1,
				},
				&ast.InfixExpression{
					Token:    token.Token{Type: token.PLUS, Literal: "+"},
					Operator: "+",
					Left: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "2"},
						Value: 2,
					},
					Right: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "2"},
						Value: 2,
					},
				},
			},
		},
		{"print(a)", "print", []string{"a"}},
		{"fn(x) { x }(x)", "fn(x) x", []string{"x"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		call, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Errorf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
		}

		functionString := call.Function.String()
		if functionString != tt.expectedFunction {
			t.Errorf("call.Function is not %s. got=%s", tt.expectedFunction, functionString)
		}

		switch expectedArgs := tt.expectedArguments.(type) {
		case []string:
			if len(call.Arguments) != len(expectedArgs) {
				t.Errorf("call does not contain %d arguments. got=%d\n", len(expectedArgs), len(call.Arguments))
			}
			for i, arg := range expectedArgs {
				testLiteralExpression(t, call.Arguments[i], arg)
			}
		case []ast.Expression:
			if len(call.Arguments) != len(expectedArgs) {
				t.Errorf("call does not contain %d arguments. got=%d\n", len(expectedArgs), len(call.Arguments))
			}
			for i, arg := range expectedArgs {
				switch exp := arg.(type) {
				case *ast.Identifier:
					testLiteralExpression(t, call.Arguments[i], exp.Value)
				case *ast.IntegerLiteral:
					testLiteralExpression(t, call.Arguments[i], exp.Value)
				case *ast.InfixExpression:
					left := extractExpressionValue(t, exp.Left)
					right := extractExpressionValue(t, exp.Right)
					testInfixExpression(t, call.Arguments[i], left, exp.Operator, right)
				}
			}
		case []any:
			for i, arg := range expectedArgs {
				switch expectedValue := arg.(type) {
				case string:
					testLiteralExpression(t, call.Arguments[i], expectedValue)
				case int:
					testLiteralExpression(t, call.Arguments[i], expectedValue)
				case int64:
					testLiteralExpression(t, call.Arguments[i], expectedValue)
				default:
					t.Fatalf("unsupported argument type for this test when using []any arguments. only identifiers or integer literals are supported")
				}
			}
		default:
			t.Fatalf("unsupported argument type for this test")
		}

	}
}

func extractExpressionValue(t *testing.T, expression ast.Expression) any {
	switch exp := expression.(type) {
	case *ast.Identifier:
		return exp.Value
	case *ast.IntegerLiteral:
		return exp.Value
	default:
		t.Fatalf("cannot extract raw value from type %T", exp)
	}
	return nil
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		expected any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. expected=%d got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator not %s. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.expected) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foo != bar;", "foo", "!=", "bar"},
		{"true != false;", true, "!=", false},
		{"true == true;", true, "==", true},
		{"false == false;", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. expected=%d got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	infixTests := []struct {
		input    string
		expected string
	}{
		{
			"5 + 5;",
			"(5 + 5)",
		},
		{
			"-a * b;",
			"((-a) * b)",
		},
		{
			"!-a;",
			"(!(-a))",
		},
		{
			"a + b + c;",
			"((a + b) + c)",
		},
		{
			"a + b - c;",
			"((a + b) - c)",
		},
		{
			"a * b * c;",
			"((a * b) * c)",
		},
		{
			"a * b / c;",
			"((a * b) / c)",
		},
		{
			"a + b / c;",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 > 4 != 3 < 4",
			"((5 > 4) != (3 < 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integer, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp is not ast.IntegerLiteral. got=%T", exp)
		return false
	}
	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d. got=%s", value, integer.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp is not ast.Boolean. got=%T", exp)
		return false
	}
	if boolean.Value != value {
		t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
		return false
	}
	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral not %t. got=%s", value, boolean.TokenLiteral())
		return false
	}
	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp is not ast.StringLiteral. got=%T", exp)
		return false
	}
	if str.Value != value {
		t.Errorf("str.Value not %s. got=%s", value, str.Value)
		return false
	}
	if str.TokenLiteral() != value {
		t.Errorf("str.TokenLiteral not %s. got=%s", value, str.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is not ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}

	t.Errorf("type of expression not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("expression is not ast.InfixExpression. got=%T", opExp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator not %s. got=%s", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
