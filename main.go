package main

import "fmt"

type Expr interface {
	isExpr()
}

// 環境
type Env map[string]int
type FuncEnv map[string]Func

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

func (Number) isExpr() {}

// シーケンス
type Seq struct {
	exprs []Expr
}

func (Seq) isExpr() {}

// 変数・代入関係
type Ident struct {
	Name string
}

func (Ident) isExpr() {}

type Assign struct {
	Name  string
	Value Expr
}

func (Assign) isExpr() {}

// 比較式
type If struct {
	Cond, Then, Else Expr
}

func (If) isExpr() {}

// 繰り返し
type While struct {
	Cond, Body Expr
}

func (While) isExpr() {}

// 関数定義
type Func struct {
	Name   string
	Params []string
	Body   Expr
}

// 関数呼び出し
type Call struct {
	Name string
	Args []Expr
}

func (Call) isExpr() {}

// プログラム
type Program struct {
	Funcs []Func
	Body  Expr
}

func Eval(expr Expr, env Env, fenv FuncEnv) (int, error) {
	switch e := expr.(type) {
	// 数値はそのまま返す
	case Number:
		return e.Value, nil
	case BinExpr:
		switch e.Op {
		// 数式はEvalMathExprで評価する
		case Add, Sub, Mul, Div:
			return EvalMathExpr(e, env, fenv)
		// 比較式はEvalCompExprで評価する
		case Lt, Gt, LtEq, GtEq, Eq, Neq:
			return EvalCompExpr(e, env, fenv)
		default:
			return 0, fmt.Errorf("unknown expression type: %T", e.Op)
		}
	case Seq:
		if len(e.exprs) == 0 {
			return 0, nil // 空のシーケンスは0を返す
		}
		var result int
		for _, ex := range e.exprs {
			v, err := Eval(ex, env, fenv)
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
		val, err := Eval(e.Value, env, fenv)
		if err != nil {
			return 0, err
		}
		env[e.Name] = val
		return val, nil
	case If:
		condVal, err := Eval(e.Cond, env, fenv)
		if err != nil {
			return 0, err
		}
		if condVal != 0 {
			return Eval(e.Then, env, fenv)
		} else {
			return Eval(e.Else, env, fenv)
		}
	case While:
		for {
			condVal, err := Eval(e.Cond, env, fenv)
			if err != nil {
				return 0, err
			}
			if condVal == 0 {
				break
			}
			_, err = Eval(e.Body, env, fenv)
			if err != nil {
				return 0, err
			}
		}
		return 0, nil // whileの値は常に0とする
	case Call:
		fn, ok := fenv[e.Name]
		if !ok {
			return 0, fmt.Errorf("undefined function: %s", e.Name)
		}
		if len(e.Args) != len(fn.Params) {
			return 0, fmt.Errorf("function %s expects %d args, got %d", e.Name, len(fn.Params), len(e.Args))
		}
		// 引数を評価
		localEnv := Env{}
		for i, arg := range e.Args {
			val, err := Eval(arg, env, fenv)
			if err != nil {
				return 0, err
			}
			localEnv[fn.Params[i]] = val
		}
		// 関数本体を新しいローカル環境で評価
		return Eval(fn.Body, localEnv, fenv)
	default:
		return 0, fmt.Errorf("unknown expression type: %T", expr)
	}
}

// 数式
func EvalMathExpr(expr BinExpr, env Env, fenv FuncEnv) (int, error) {
	l, err := Eval(expr.Left, env, fenv)
	if err != nil {
		return 0, err
	}
	r, err := Eval(expr.Right, env, fenv)
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
func EvalCompExpr(expr BinExpr, env Env, fenv FuncEnv) (int, error) {
	l, err := Eval(expr.Left, env, fenv)
	if err != nil {
		return 0, err
	}
	r, err := Eval(expr.Right, env, fenv)
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

func EvalProgram(p Program) (int, error) {
	env := Env{}
	fenv := FuncEnv{}

	for _, fn := range p.Funcs {
		fenv[fn.Name] = fn
	}

	return Eval(p.Body, env, fenv)
}

func main() {
	env := make(Env)
	fenv := make(FuncEnv)

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
	result, err := Eval(expr, env, fenv)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// (1 + 2) * 3 > 10 を表す式を構築
	expr2 := BinExpr{
		Op: Gt,
		Left: BinExpr{
			Op: Mul,
			Left: BinExpr{
				Op:    Add,
				Left:  Number{Value: 1},
				Right: Number{Value: 2},
			},
			Right: Number{Value: 3},
		},
		Right: Number{Value: 10},
	}
	result2, err := Eval(expr2, env, fenv)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result2)
	}

	// シーケンスの例: (1 + 2; 3 * 4) を表す式を構築
	expr3 := Seq{
		exprs: []Expr{
			BinExpr{
				Op:    Add,
				Left:  Number{Value: 1},
				Right: Number{Value: 2},
			},
			BinExpr{
				Op:    Mul,
				Left:  Number{Value: 3},
				Right: Number{Value: 4},
			},
		},
	}
	result3, err := Eval(expr3, env, fenv)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result3)
	}

	// 変数と代入の例: x = 5; x + 3 を表す式を構築
	expr4 := Seq{
		exprs: []Expr{
			Assign{
				Name:  "x",
				Value: Number{Value: 5},
			},
			BinExpr{
				Op:    Add,
				Left:  Ident{Name: "x"},
				Right: Number{Value: 3},
			},
		},
	}
	result4, err := Eval(expr4, env, fenv)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result4)
	}

	// 比較式の例: if (x > 3) then 1 else 0 を表す式を構築
	expr5 := If{
		Cond: BinExpr{
			Op:    Gt,
			Left:  Ident{Name: "x"},
			Right: Number{Value: 3},
		},
		Then: Number{Value: 1},
		Else: Number{Value: 0},
	}
	result5, err := Eval(expr5, env, fenv)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result5)
	}

	// 繰り返しの例: while (x < 10) do x = x + 1 を表す式を構築
	expr6 := Seq{
		exprs: []Expr{
			Assign{
				Name:  "x",
				Value: Number{Value: 0},
			},
			While{
				Cond: BinExpr{
					Op:    Lt,
					Left:  Ident{Name: "x"},
					Right: Number{Value: 10},
				},
				Body: Assign{
					Name: "x",
					Value: BinExpr{
						Op:    Add,
						Left:  Ident{Name: "x"},
						Right: Number{Value: 1},
					},
				},
			},
			Ident{Name: "x"}, // 最終的な値を返すために変数xを評価
		},
	}
	result6, err := Eval(expr6, env, fenv)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result6)
	}

	// 関数定義と呼び出しの例: func add(a, b) { a + b }; add(3, 4) => 7
	prog := Program{
		Funcs: []Func{
			{
				Name:   "add",
				Params: []string{"a", "b"},
				Body: BinExpr{
					Op:    Add,
					Left:  Ident{Name: "a"},
					Right: Ident{Name: "b"},
				},
			},
		},
		Body: Call{Name: "add", Args: []Expr{Number{Value: 3}, Number{Value: 4}}},
	}
	result7, err := EvalProgram(prog)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result7)
	}

	// 再帰関数の例: func fact(n) { if (n <= 1) 1 else n * fact(n - 1) }; fact(5) => 120
	prog2 := Program{
		Funcs: []Func{
			{
				Name:   "fact",
				Params: []string{"n"},
				Body: If{
					Cond: BinExpr{Op: LtEq, Left: Ident{Name: "n"}, Right: Number{Value: 1}},
					Then: Number{Value: 1},
					Else: BinExpr{
						Op:   Mul,
						Left: Ident{Name: "n"},
						Right: Call{
							Name: "fact",
							Args: []Expr{
								BinExpr{Op: Sub, Left: Ident{Name: "n"}, Right: Number{Value: 1}},
							},
						},
					},
				},
			},
		},
		Body: Call{Name: "fact", Args: []Expr{Number{Value: 5}}},
	}
	result8, err := EvalProgram(prog2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result8)
	}
}
