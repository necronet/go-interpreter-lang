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

func TestStringLiteral(t *testing.T) {
    input := `"Hello world!"`
    evaluated := testEval(input)

    str, ok := evaluated.(*object.String)
    if !ok {
        t.Fatalf("object is not string. got=%T (%+v)", evaluated, evaluated)
    }
    if str.Value != "Hello world!" {
        t.Errorf("String has the wrong value. got=%q but expected %q", str.Value, input)
    }
}

func TestStringConcatentation(t *testing.T) {
    input := `"Hello" + " " + "World!"`
    
    evaluated := testEval(input)

    str, ok := evaluated.(*object.String)
    if !ok {
        t.Fatalf("object is not string. got=%T (%+v)", evaluated, evaluated)
    }
    if str.Value != "Hello World!" {
        t.Errorf("String has the wrong value. got=%q but expected %q", str.Value, input)
    }
}

func TestLetStatements(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"let a = 5; a", 5},
        {"let a = 5 * 5; a;", 25},
        {"let a = 5; let b = a; b;", 5},
        {"let a = 5; let b = a; let c = a + b + 5; c;", 15},
    }

    for _, tt := range tests{
        testIntegerObject(t, testEval(tt.input), tt.expected)
    }
}

func TestFunctionObject(t *testing.T) {
    input := "fn(x) { x + 2; } ;"

    evaluated := testEval(input)

    fn, ok := evaluated.(*object.Function)
    if !ok {
        t.Fatalf("object is not a Function. got=%T (%+v)", evaluated, evaluated)
    }

    if len(fn.Parameters) != 1 {
        t.Fatalf("function has wrong parameters. Parameters=+%v", fn.Parameters)
    }

    if fn.Parameters[0].String() != "x" {
        t.Fatalf("parameter is not 'x' as expected. got=%q", fn.Parameters[0])
    }

    expectedBody := "(x + 2)"

    if fn.Body.String() != expectedBody {
        t.Fatalf("body is not %q. got=%q instead", expectedBody, fn.Body.String())
    }
}

func TestFunctionApplication(t *testing.T) {

    tests := []struct{
        input string
        expected int64
    }{
    { "let identity = fn(x) { x; }; identity(5);", 5},
    { "let identity = fn(x) { x; return x; }; identity(5);", 5},
    { "let double = fn(x) { x*2; }; double(5);", 10},
    { "let add = fn(x, y) { x+y; }; add(3, 5);", 8},
    { "let add = fn(x, y) { x+y; }; add(5+5, add(5,5));", 20},
    {"fn(x) { x; }(5)", 5},
    }

    for _, tt := range tests {
        testIntegerObject(t, testEval(tt.input), tt.expected)
    }
}

func TestErrorHandling(t *testing.T) {
    tests := []struct {
        input string
        expectedMessage string
    }{
        { "5 + true", "type mismatch: INTEGER + BOOLEAN"},
        { "5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
        { "-true", "unknown operator: -BOOLEAN"},
        { "true + false", "unknown operator: BOOLEAN + BOOLEAN"},
        { "5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
        { "if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
        {
            ` if (10 > 1) {
                if (10 >1 ) {
                    return true + false;
                }
            }
            return 1;
            `, "unknown operator: BOOLEAN + BOOLEAN",
        },
        { "foobar", "identifier not found: foobar"},
        {`"Hello" - "World"`, "unknown operator: STRING - STRING"},
    }

    for _, tt := range tests {

        evaluated := testEval(tt.input)

        errObj, ok := evaluated.(*object.Error)
        if !ok {
            t.Errorf("no error object returned. got=%T(%+v)",
            evaluated, evaluated)
            continue
        }

        if errObj.Message != tt.expectedMessage {
            t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
        }

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
    env := object.NewEnvironment()

    return Eval(program, env)
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
