package main

import (
	"encoding/json"
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
			return eval(n.X) + eval(n.Y)
		case token.SUB:
			return eval(n.X) - eval(n.Y)
		case token.MUL:
			return eval(n.X) * eval(n.Y)
		case token.QUO:
			return eval(n.X) / eval(n.Y)
		}
	case *ast.UnaryExpr:
		if n.Op == token.SUB {
			return -eval(n.X)
		}
	case *ast.BasicLit:
		if n.Kind == token.INT || n.Kind == token.FLOAT {
			val, err := strconv.Atoi(n.Value)
			if err != nil {
				fmt.Println(err) // TODO: We should bubble this error up
			}
			return val
		}
	case *ast.ParenExpr:
		return eval(n.X)
	}
	return 0 // TODO: Unhandled type: should return error
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
		fmt.Println("Handled error:", err) // TODO: Remove after debugging
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}
	ast.Print(fileSet, tree)
	fmt.Println("Result:", eval(tree))
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
