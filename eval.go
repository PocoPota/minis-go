package main

import "fmt"

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
		if len(e.Exprs) == 0 {
			return 0, nil // 空のシーケンスは0を返す
		}
		var result int
		for _, ex := range e.Exprs {
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