package main

import "fmt"

type Expr interface {
	isExpr()
}

// Opの定義
type Op int

const (
	Add Op = iota
	Sub
	Mul
	Div
)

type BinExpr struct {
	Op
	Left, Right Expr
}

type Number struct {
	Value int
}

func (BinExpr) isExpr() {}
func (Number) isExpr()  {}

func Eval(expr Expr) (int, error) {
	switch e := expr.(type) {
	case Number:
		return e.Value, nil
	case BinExpr:
		return EvalMathExpr(e)
	default:
		return 0, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func EvalMathExpr(expr BinExpr) (int, error) {
	l, err := Eval(expr.Left)
	if err != nil {
		return 0, err
	}
	r, err := Eval(expr.Right)
	if err != nil {
		return 0, err
	}

	switch expr.Op {
	case Add:
		return l + r, nil
	case Sub:
		return l - r, nil
	case Mul:
		return l * r, nil
	case Div:
		if r == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return l / r, nil
	default:
		return 0, fmt.Errorf("unknown operator: %v", expr.Op)
	}
}

func main() {
	// (1 + 2 * 3) を表す式を構築
	expr := BinExpr{
		Op:   Add,
		Left: Number{Value: 1},
		Right: BinExpr{
			Op:    Mul,
			Left:  Number{Value: 2},
			Right: Number{Value: 3},
		},
	}
	result, err := Eval(expr)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
