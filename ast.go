package main

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
	Exprs []Expr
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

// 条件分岐
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