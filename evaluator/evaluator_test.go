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

func TestBuiltinFunctions(t *testing.T) {

    tests := []struct {
        input string
        expected interface {}
    }{
        {`len("")`, 0},
        {`len("four")`, 4},
        {`len("hello world")`, 11},
        {`len(1)`, "argument to `len` not supported, got INTEGER"},
        {`len("four", "trois")`, "wrong number of arguments. got=2, want=1"},
    }

    for _, tt := range tests{
        evaluated := testEval(tt.input)

        switch expected := tt.expected.(type) {

        case int:
            testIntegerObject(t, evaluated, int64(expected))
        case string:
            errObj, ok := evaluated.(*object.Error)
            if !ok {
                t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
                continue
            }
            if errObj.Message != expected {
                t.Errorf("Message is wrong. expected=%q, got=%q", expected, errObj.Message)
            }
        }

    }
}

func TestArrayListerals(t *testing.T) {
    input := "[1, 2, 3, 2*2]"

    evaluated := testEval(input)
    result, ok := evaluated.(*object.Array)
    if !ok {
        t.Fatalf("object is not Array. got=%t (%+v)", evaluated, evaluated)
    }

    if len(result.Elements) != 4 {
        t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
    }
    testIntegerObject(t, result.Elements[0], 1)
    testIntegerObject(t, result.Elements[1], 2)
    testIntegerObject(t, result.Elements[2], 3)
    testIntegerObject(t, result.Elements[3], 4)
}

func TestArrayIndexExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {
            "[1, 2, 3][0]",
            1,
        },
        {
            "[1, 2, 3][1]",
            2,
        },
        {
            "[1, 2, 3][2]",
            3,
        },
        {
            "let i = 0; [1][i];",
            1,
        },
        {
            "[1, 2, 3][1 + 1];",
            3,
        },
        {
            "let myArray = [1, 2, 3]; myArray[2];",
            3,
        },
        {
            "let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
            6,
        },
        {
            "let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
            2,
        },
        {
            "[1, 2, 3][3]",
            nil,
        },
        {
            "[1, 2, 3][-1]",
            nil,
        },
    }
for _, tt := range tests {
    evaluated := testEval(tt.input)
    integer, ok := tt.expected.(int)
    if ok {
        testIntegerObject(t, evaluated, int64(integer))
    } else {
        testNullObject(t, evaluated)
    }
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
