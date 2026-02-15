package main

import "fmt"

// 式
type Expr interface {
	isExpr()
}

// 数値
type Number struct {
	Value int
}

// 加算
type Add struct {
	Left  Expr
	Right Expr
}

// 減算
type Sub struct {
	Left  Expr
	Right Expr
}

// 乗算
type Mul struct {
	Left  Expr
	Right Expr
}

// 除算
type Div struct {
	Left  Expr
	Right Expr
}

// isExprメソッドの追加
func (Number) isExpr() {}
func (Add) isExpr()    {}
func (Sub) isExpr()    {}
func (Mul) isExpr()    {}
func (Div) isExpr()    {}

func Eval(expr Expr) (int, error){
	switch e := expr.(type) {
	
	case Number:
		return e.Value, nil
	
	case Add:
		l, err := Eval(e.Left)
		if err != nil {
			return 0, err
		}
		r, err := Eval(e.Right)
		if err != nil {
			return 0, err
		}
		return l + r, nil
			
	case Sub:
		l, err := Eval(e.Left)
		if err != nil {
			return 0, err
		}
		r, err := Eval(e.Right)
		if err != nil {
			return 0, err
		}
		return l - r, nil
			
	case Mul:
		l, err := Eval(e.Left)
		if err != nil {
			return 0, err
		}
		r, err := Eval(e.Right)
		if err != nil {
			return 0, err
		}
		return l * r, nil
			
	case Div:
		l, err := Eval(e.Left)
		if err != nil {
			return 0, err
		}
		r, err := Eval(e.Right)
		if err != nil {
			return 0, err
		}
		if r == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return l / r, nil
		
	default:
		return 0, fmt.Errorf("unknown expression")
	}
}

func main() {
	// (3 + 5) * (10 - 2) = 64
	expr := Mul{
		Left: Add{
			Left: Number{Value: 3},
			Right: Number{Value: 5},
		},
		Right: Sub{
			Left: Number{Value: 10},
			Right: Number{Value: 2},
		},
	}
	result, err := Eval(expr)
	if err != nil{
		fmt.Println("Error:", err)
	}

	fmt.Println("Result:", result)
}
