package main

import "fmt"

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
		Exprs: []Expr{
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
		Exprs: []Expr{
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
		Exprs: []Expr{
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
