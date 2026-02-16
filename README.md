# Minis Go
Minis Go は [MiniS](https://github.com/kmizu/minis/tree/main) を参考に Go で実装された問プログラミング言語です。

## 対応する構文

### 数値
```go
Number{Value: 3} // 3
```

### 四則演算
```go
// 1 + 2 => 3
BinExpr{
  Op: Add,
  Left: Number{Value: 1},
  Right: Number{Value: 2}
}
```

`Op` では `Add`, `Sub`, `Mul`, `Div` を指定することができます。

### 比較式
```go
// 1 > 2 => 0 (false)
BinExpr{
  Op: Lt,
  Left: Number{Value: 1},
  Right: Number{Value: 2}
}
```

`Op` では `Lt` (<), `Gt` (>), `LtEq` (<=), `GtEq` (>=), `Eq` (==), `Neq` (!=) を指定することが出来ます。

`0` (false), `1` (true) が返されます。

### 代入式
```go
// x = 5
Assign{
  Name: "x",
  Value: Number{Value: 5}
}

// x を使う
Ident{Name: "x"}
```

### 連接式
```go
// x = 5; x + 3
Seq{
  Exprs: []Expr{
    Assign{
      Name: "x",
      Value: Number{Value: 5}
    },
    BinExpr{
      Op: Add,
      Left: Ident{Name: "x"},
      Right: Number{Value: 3}
    }
  }
}
```

### 条件分岐式
```go
// if (x > 3) then 1 else 0
If{
  Cond: BinExpr{
    Op: Gt,
    Left: Ident{Name: "x"},
    Right: Number{Value: 3}
  },
  Then: Number{Value: 1},
  Else: Number{Value: 0}
}
```

### 繰り返し式
```go
// x = 0; while(x < 10) {x = x + 1}; x
Seq{
  Exprs: []Expr{
    Assign: {Name: "x", Value: Number{Value: 0}}
  },
  While{
    Cond: BinExpr{
      Op: Lt,
      Left: Ident{Name: "x"},
      Right: Number{Value: 10}
    },
    Body: Assign{
      Name: "x",
      Value: BinExpr{
        Op: Add,
        Left: Ident{Name: "x"},
        Right: Number{Value: 1}
      }
    }
  },
  Ident{Name: "x"}
}
```

### 関数
```go
// func add(a, b) { a + b }; add(3, 4)
Program{
	Funcs: []Func{
		{
			Name: "add",
			Params: []string{"a", "b"},
			Body: BinExpr{
				Op: Add,
				Left: Ident{Name: "a"},
				Right: Ident{Name: "b"},
			},
		},
	},
	Body: Call{Name: "add", Args: []Expr{Number{Value: 3}, Number{Value: 4}}},
}
```