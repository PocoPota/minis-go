package main

import (
	"encoding/json"
	"fmt"
)

func ParseProgram(data []byte) (Program, error) {
	var raw struct {
		Funcs []json.RawMessage `json:"funcs"`
		Body  json.RawMessage   `json:"body"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Program{}, fmt.Errorf("parse program: %w", err)
	}

	funcs := make([]Func, len(raw.Funcs))
	for i, rf := range raw.Funcs {
		f, err := parseFunc(rf)
		if err != nil {
			return Program{}, err
		}
		funcs[i] = f
	}

	body, err := parseExpr(raw.Body)
	if err != nil {
		return Program{}, err
	}

	return Program{Funcs: funcs, Body: body}, nil
}

func parseFunc(data json.RawMessage) (Func, error) {
	var raw struct {
		Name   string          `json:"name"`
		Params []string        `json:"params"`
		Body   json.RawMessage `json:"body"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Func{}, fmt.Errorf("parse func: %w", err)
	}

	body, err := parseExpr(raw.Body)
	if err != nil {
		return Func{}, err
	}

	return Func{Name: raw.Name, Params: raw.Params, Body: body}, nil
}

var opMap = map[string]Op{
	"Add":  Add,
	"Sub":  Sub,
	"Mul":  Mul,
	"Div":  Div,
	"Lt":   Lt,
	"Gt":   Gt,
	"LtEq": LtEq,
	"GtEq": GtEq,
	"Eq":   Eq,
	"Neq":  Neq,
}

func parseExpr(data json.RawMessage) (Expr, error) {
	var head struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return nil, fmt.Errorf("parse expr type: %w", err)
	}

	switch head.Type {
	case "num":
		var v struct {
			Value int `json:"value"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return Number{Value: v.Value}, nil

	case "ident":
		var v struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return Ident{Name: v.Name}, nil

	case "bin":
		var v struct {
			Op    string          `json:"op"`
			Left  json.RawMessage `json:"left"`
			Right json.RawMessage `json:"right"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		op, ok := opMap[v.Op]
		if !ok {
			return nil, fmt.Errorf("unknown op: %s", v.Op)
		}
		left, err := parseExpr(v.Left)
		if err != nil {
			return nil, err
		}
		right, err := parseExpr(v.Right)
		if err != nil {
			return nil, err
		}
		return BinExpr{Op: op, Left: left, Right: right}, nil

	case "if":
		var v struct {
			Cond json.RawMessage `json:"cond"`
			Then json.RawMessage `json:"then"`
			Else json.RawMessage `json:"else"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		cond, err := parseExpr(v.Cond)
		if err != nil {
			return nil, err
		}
		then, err := parseExpr(v.Then)
		if err != nil {
			return nil, err
		}
		els, err := parseExpr(v.Else)
		if err != nil {
			return nil, err
		}
		return If{Cond: cond, Then: then, Else: els}, nil

	case "while":
		var v struct {
			Cond json.RawMessage `json:"cond"`
			Body json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		cond, err := parseExpr(v.Cond)
		if err != nil {
			return nil, err
		}
		body, err := parseExpr(v.Body)
		if err != nil {
			return nil, err
		}
		return While{Cond: cond, Body: body}, nil

	case "assign":
		var v struct {
			Name  string          `json:"name"`
			Value json.RawMessage `json:"value"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		val, err := parseExpr(v.Value)
		if err != nil {
			return nil, err
		}
		return Assign{Name: v.Name, Value: val}, nil

	case "seq":
		var v struct {
			Exprs []json.RawMessage `json:"exprs"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		exprs := make([]Expr, len(v.Exprs))
		for i, raw := range v.Exprs {
			e, err := parseExpr(raw)
			if err != nil {
				return nil, err
			}
			exprs[i] = e
		}
		return Seq{Exprs: exprs}, nil

	case "call":
		var v struct {
			Name string            `json:"name"`
			Args []json.RawMessage `json:"args"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		args := make([]Expr, len(v.Args))
		for i, raw := range v.Args {
			a, err := parseExpr(raw)
			if err != nil {
				return nil, err
			}
			args[i] = a
		}
		return Call{Name: v.Name, Args: args}, nil

	default:
		return nil, fmt.Errorf("unknown expr type: %s", head.Type)
	}
}
