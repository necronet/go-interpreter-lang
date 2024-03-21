package evaluator

import(
    "necronet.info/interpreter/lexer"
    "necronet.info/interpreter/object"
    "necronet.info/interpreter/parser"
    "testing"
)

func TestEvalIntegerExpression(t *testing.T) {
    tests :=[]struct {
        input string
        expected int64
    }{
        {"5", 5},
        {"10", 10},
        {"-5", -5},
        {"-10", -10},
        {"5+5+5+5",20},
        {"5+5+5-5",10},
        {"3*10",30},
        {"15/5",3},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}

func TestReturnStatements(t *testing.T) {
    tests := []struct{
        input string
        expected int64
    }{
        {"return 10;", 10},
        {"return 10; 9;", 10},
        {"return 2 * 5;", 10},
        {"9; return 10; 2*5", 10},
    }
    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}

func TestIfElseExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"if (true) { 9 }", 9},
        {"if (false) { 11 }", nil},
        {"if (1) { 12 }", 12},
        {"if (1 < 2) { 10 }", 10},
        {"if (1 > 2) { 13 }", nil},
    }

    for _, tt := range tests {

        evaluated := testEval(tt.input)
        integer, ok := tt.expected.(int)

        if ok {
            testIntegerObject(t, evaluated, int64(integer))
        }else {
            testNullObject(t, evaluated)
        }
    }
}

func testNullObject(t *testing.T, obj object.Object) bool {

    if obj != NULL {
        t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
        return false
    }
    return true
}

func TestEvalBooleanExpression(t *testing.T){
    tests := []struct {
        input string
        expected bool
    }{
        {"true", true},{"false",false},
        {" 1 < 2", true},
        {" 1 > 2", false},
        {" 1 < 1", false},
        {" 1 == 2", false},
        {" 1 != 2", true},
        {"true == true", true},
        {"false == true", false},
        {"true == false", false},
        {"true != false", true},
    }

    for _, tt  := range tests {
        evaluated := testEval(tt.input)
        testBooleanObject(t, evaluated, tt.expected)
    }
}
func testEval(input string) object.Object {

    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()

    return Eval(program)
}

func testBangOperator(t *testing.T) {

    tests := []struct{ 
        input string
        expected bool
    }{
        {"!true", false},
        {"!false", true},
        {"!5", false},
        {"!!true", true},
        {"!false", true},
        {"!!5", true},
    }

    for _, tt := range tests {

        evaluated := testEval(tt.input)
        testBooleanObject(t, evaluated, tt.expected)
    }
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
    result, ok := obj.(*object.Integer)

    if !ok {
        t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
        return false
    }
    if result.Value != expected {
        t.Errorf("Object has wrong value. got=%d, want %d", result.Value, expected)
        return false
    }

    return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
result, ok := obj.(*object.Boolean)

if !ok {
    t.Errorf("object is not boolean. got=%T (%+v)", obj, obj)
    return false
}
if result.Value != expected {
    t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
    return false
}
return true

}
