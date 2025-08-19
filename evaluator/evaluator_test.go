package evaluator

import (
	"fmt"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"true == false", false},
		{"true != true", false},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"(1 < 2) == (2 > 1)", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world"`, "hello world"},
		{`"hello \"world\""`, "hello \"world\""},
		{`"hello\nworld"`, "hello\nworld"},
		{`"hello\tworld"`, "hello\tworld"},
		{`"hello" + " " + "world"`, "hello world"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestArrayLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{
		{"[1]", []any{1}},
		{"[1, 2]", []any{1, 2}},
		{"[1, 2, \"foo\"]", []any{1, 2, "foo"}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"[1][0]", 1},
		{"[1, 2][1]", 2},
		{"[1, 2, \"foo\"][2]", "foo"},
		{"let myArray = [1, 2, \"foo\"]; myArray[2];", "foo"},
		{"[1, 2, 3][3]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testStringObject(t, evaluated, expected)
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!0", true},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!!0", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			fmt.Printf("\tinput: %s\n", tt.input)
		}
	}
}

func TestIfExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 5 }", 5},
		{"if (false) { 5 } else { 10 }", 10},
		{"if (false) { 5 }", nil},
		{"if (1) { 1 } else { 0 }", 1},
		{"if (0) { 1 } else { 0 }", 0},
		{"if (1 < 2) { 1 } else { 2 }", 1},
		{"if (1 > 2) { 1 } else { 2 }", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		expected, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(expected))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, int64(tt.expected))
	}
}

func TestAssignmentStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let a = 5; a = 10; a", 10},
		{"let a = 5 * 5; a = a + a; a", 50},
		{"let a = 5; let b = 0; b = a; b", 5},
		{"let a = 5; let b = a; let c = 0; c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, int64(tt.expected))
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`if (10 > 1) {
        if (10 > 1) {
          return 10
        }
        return 1
      }
      `,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not a Function, got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("fn does not have %d parameters, got=%d", 1, len(fn.Parameters))
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("fn.Parameters[0] is not named %q, got=%q", "x", fn.Parameters[0].String())
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("fn.Body is not %q, got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { return x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments, got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error, got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != expected {
				t.Errorf("wrong error message, expeted=%q, got=%q", expected, errObj.Message)
			}

		}
	}
}

func TestErrorHandlig(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true;",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (true + false) { true; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
        if (10 > 1) {
          true + false
          return 10 
        }
        return 1
      }
      `,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier foobar is undefined",
		},
		{
			"foo = 1",
			"assign to an undefined identifier foo",
		},
		{
			`"hello" - "world"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error returned, got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q got=%q", tt.expected, errObj.Message)
		}
	}
}

func testEval(input string) object.Object {
	env := object.NewEnvironment()
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%d, want=%d ", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%t, want=%t ", result.Value, expected)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%s, want=%s ", result.Value, expected)
		return false
	}

	return true
}

func testArrayObject(t *testing.T, obj object.Object, expected []any) bool {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not Array, got=%T (%+v)", obj, obj)
		return false
	}

	if len(result.Elements) != len(expected) {
		t.Errorf("array has wrong number of elements, got=%d, want=%d ", len(result.Elements), len(expected))
		return false
	}

	for i, el := range result.Elements {
		switch expected := expected[i].(type) {
		case int:
			testIntegerObject(t, el, int64(expected))
		case string:
			testStringObject(t, el, expected)
		}
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	result, ok := obj.(*object.Null)
	if !ok {
		t.Errorf("object is not Null, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Type() != object.NullType {
		t.Errorf("object.Type is not object.NullObj, got=%s", obj.Type())
		return false
	}
	return true
}
