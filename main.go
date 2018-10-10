package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"strconv"
	"strings"
)

const apiRoot = "/api/"

type CalcRequest struct {
	Expression string `json:"expression"`
}

func eval(node ast.Node) int {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		switch n.Op {
		case token.ADD:
			xVal, err := eval(n.X)
			if err != nil {
				return 0, err
			}
			yVal, err := eval(n.Y)
			if err != nil {
				return 0, err
			}
			return xVal + yVal, nil
		case token.SUB:
			xVal, err := eval(n.X)
			if err != nil {
				return 0, err
			}
			yVal, err := eval(n.Y)
			if err != nil {
				return 0, err
			}
			return xVal - yVal, nil
		case token.MUL:
			xVal, err := eval(n.X)
			if err != nil {
				return 0, err
			}
			yVal, err := eval(n.Y)
			if err != nil {
				return 0, err
			}
			return xVal * yVal, nil
		case token.QUO:
			xVal, err := eval(n.X)
			if err != nil {
				return 0, err
			}
			yVal, err := eval(n.Y)
			if err != nil {
				return 0, err
			}
			return xVal / yVal, nil // TODO: Currently using integer division. Maybe change it?
		}
	case *ast.UnaryExpr:
		if n.Op == token.SUB {
			val, err := eval(n.X)
			return -val, err
		}
	case *ast.BasicLit:
		if n.Kind == token.INT || n.Kind == token.FLOAT {
			val, err := strconv.Atoi(n.Value)
			return val, err
		}
	case *ast.ParenExpr:
		return eval(n.X)
	}
	return 0, errors.New("unhandled node type")
}

func handleCalc(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var request CalcRequest
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	fileSet := token.NewFileSet()
	tree, err := parser.ParseExpr(request.Expression)
	if err != nil {
		fmt.Println("Bad request, handled error:", err)
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}
	ast.Print(fileSet, tree)
	result, err := eval(tree)
	if err != nil {
		fmt.Println("Error during evaluation, probably caused by a bad request", err)
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}
	fmt.Println("Result:", result)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	urlParts := strings.Split(strings.TrimPrefix(r.URL.Path, apiRoot), "/")

	// Very simple router
	switch urlParts[0] {
	case "calc":
		if len(urlParts) > 1 && len(urlParts[1]) != 0 {
			http.Error(w, "404 not found", http.StatusNotFound) // Disallow /api/calc/*
			break
		}
		handleCalc(w, r)
		break
	case "coffee":
		if len(urlParts) > 1 && len(urlParts[1]) != 0 {
			http.Error(w, "404 not found", http.StatusNotFound) // Disallow /api/coffee/*
			break
		}
		http.Error(w, "418 I'm a teapot", http.StatusTeapot)
	default:
		http.Error(w, "404 not found", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc(apiRoot, apiHandler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
