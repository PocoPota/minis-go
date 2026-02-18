# Minis Go
Minis Go は [MiniS](https://github.com/kmizu/minis/tree/main) を参考に Go で実装されたトイプログラミング言語です。

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

## JSON でプログラムを書く

JSON ファイルにプログラムを記述し、コマンドライン引数で指定して実行できます。

```sh
go run . fact.json
```

### 式の表現

各式は `type` フィールドで種類を指定します。

| type | 説明 | フィールド |
|---|---|---|
| `num` | 数値 | `value` |
| `ident` | 変数参照 | `name` |
| `bin` | 二項演算 | `op`, `left`, `right` |
| `if` | 条件分岐 | `cond`, `then`, `else` |
| `while` | 繰り返し | `cond`, `body` |
| `assign` | 代入 | `name`, `value` |
| `seq` | 連接 | `exprs` |
| `call` | 関数呼び出し | `name`, `args` |

`bin` の `op` には `Add`, `Sub`, `Mul`, `Div`, `Lt`, `Gt`, `LtEq`, `GtEq`, `Eq`, `Neq` を指定できます。

### プログラムの構造

トップレベルは `funcs`（関数定義の配列）と `body`（メインの式）で構成されます。

### 例: 階乗（fact(5) = 120）

```json
{
  "funcs": [
    {
      "name": "fact",
      "params": ["n"],
      "body": {
        "type": "if",
        "cond": {
          "type": "bin",
          "op": "LtEq",
          "left": { "type": "ident", "name": "n" },
          "right": { "type": "num", "value": 1 }
        },
        "then": { "type": "num", "value": 1 },
        "else": {
          "type": "bin",
          "op": "Mul",
          "left": { "type": "ident", "name": "n" },
          "right": {
            "type": "call",
            "name": "fact",
            "args": [
              {
                "type": "bin",
                "op": "Sub",
                "left": { "type": "ident", "name": "n" },
                "right": { "type": "num", "value": 1 }
              }
            ]
          }
        }
      }
    }
  ],
  "body": {
    "type": "call",
    "name": "fact",
    "args": [{ "type": "num", "value": 5 }]
  }
}
```