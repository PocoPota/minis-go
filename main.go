package main

import "fmt"

type Expr interface {
	isExpr()
}

// Opの定義
type Op int

const (
	Add  Op = iota // +
	Sub            // -
	Mul            // *
	Div            // /
	Lt             // <
	Gt             // >
	LtEq           // <=
	GtEq           // >=
	Eq             // ==
	Neq            // !=
)

// 二項式
type BinExpr struct {
	Op
	Left, Right Expr
}
func (BinExpr) isExpr() {}

// 数値
type Number struct {
	Value int
}
func (Number) isExpr()  {}

// シーケンス
type Seq struct {
	exprs []Expr
}
func (Seq) isExpr() {}

// 変数・代入関係
type Env map[string]int

type Ident struct {
	Name string
}
func (Ident) isExpr() {}

type Assign struct {
	Name string
	Value Expr
}
func (Assign) isExpr() {}

func Eval(expr Expr, env Env) (int, error) {
	switch e := expr.(type) {
	// 数値はそのまま返す
	case Number:
		return e.Value, nil
	case BinExpr:
		switch e.Op {
		// 数式はEvalMathExprで評価する
		case Add, Sub, Mul, Div:
			return EvalMathExpr(e, env)
		// 比較式はEvalCompExprで評価する
		case Lt, Gt, LtEq, GtEq, Eq, Neq:
			return EvalCompExpr(e, env)
		default:
			return 0, fmt.Errorf("unknown expression type: %T", e.Op)
		}
	case Seq:
		if len(e.exprs) == 0 {
			return 0, nil // 空のシーケンスは0を返す
		}
		var result int
		for _, ex := range e.exprs{
			v, err := Eval(ex, env)
			if err != nil {
				return 0, err
			}
			result = v // 最後の式の値を返す
		}
		return result, nil
	case Ident:
		val, err := env[e.Name]
		if !err {
			return 0, fmt.Errorf("undefined ident: %s", e.Name)
		}
		return val, nil
	case Assign:
		val, err := Eval(e.Value, env)
		if err != nil {
			return 0, err
		}
		env[e.Name] = val
		return val, nil
	default:
		return 0, fmt.Errorf("unknown expression type: %T", expr)
	}
}

// 数式
func EvalMathExpr(expr BinExpr, env Env) (int, error) {
	l, err := Eval(expr.Left, env)
	if err != nil {
		return 0, err
	}
	r, err := Eval(expr.Right, env)
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

// 比較式 (0: false, 1: true)
func EvalCompExpr(expr BinExpr, env Env) (int, error) {
	l, err := Eval(expr.Left, env)
	if err != nil {
		return 0, err
	}
	r, err := Eval(expr.Right, env)
	if err != nil {
		return 0, err
	}

	switch expr.Op {
	case Lt:
		if l < r {
			return 1, nil
		}
		return 0, nil
	case Gt:
		if l > r {
			return 1, nil
		}
		return 0, nil
	case LtEq:
		if l <= r {
			return 1, nil
		}
		return 0, nil
	case GtEq:
		if l >= r {
			return 1, nil
		}
		return 0, nil
	case Eq:
		if l == r {
			return 1, nil
		}
		return 0, nil
	case Neq:
		if l != r {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unknown operator: %v", expr.Op)
	}
}

func main() {
	env := make(Env)

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
	result, err := Eval(expr, env)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// (1 + 2) * 3 > 10 を表す式を構築
	expr2 := BinExpr{
		Op:   Gt,
		Left: BinExpr{
			Op:   Mul,
			Left: BinExpr{
				Op:   Add,
				Left: Number{Value: 1},
				Right: Number{Value: 2},
			},
			Right: Number{Value: 3},
		},
		Right: Number{Value: 10},
	}
	result2, err := Eval(expr2, env)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result2)
	}

	// シーケンスの例: (1 + 2; 3 * 4) を表す式を構築
	expr3 := Seq{
		exprs: []Expr{
			BinExpr{
				Op:   Add,
				Left: Number{Value: 1},
				Right: Number{Value: 2},
			},
			BinExpr{
				Op:   Mul,
				Left: Number{Value: 3},
				Right: Number{Value: 4},
			},
		},
	}
	result3, err := Eval(expr3, env)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result3)
	}

	// 変数と代入の例: x = 5; x + 3 を表す式を構築
	expr4 := Seq{
		exprs: []Expr{
			Assign{
				Name: "x",
				Value: Number{Value: 5},
			},
			BinExpr{
				Op:   Add,
				Left: Ident{Name: "x"},
				Right: Number{Value: 3},
			},
		},
	}
	result4, err := Eval(expr4, env)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result4)
	}
}
