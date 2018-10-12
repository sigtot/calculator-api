package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sigtot/rest-calculator/file_handling"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const apiRoot = "/api/"

const historiesFileName = "histories.db"
const defaultNumHistoryLines = 5

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	Result float64 `json:"result"`
}

type Calculation struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

// Eval can evaluate mathematical expressions with the tokens +, -, *, / and ()
func Eval(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		var xVal float64
		var yVal float64
		var errX error
		var errY error
		switch n.Op {
		case token.ADD:
			xVal, errX = Eval(n.X)
			yVal, errY = Eval(n.Y)
			return xVal + yVal, orError(errX, errY)
		case token.SUB:
			xVal, errX = Eval(n.X)
			yVal, errY = Eval(n.Y)
			return xVal - yVal, orError(errX, errY)
		case token.MUL:
			xVal, errX = Eval(n.X)
			yVal, errY = Eval(n.Y)
			return xVal * yVal, orError(errX, errY)
		case token.QUO:
			xVal, errX = Eval(n.X)
			yVal, errY = Eval(n.Y)
			if yVal != 0 {
				return xVal / yVal, orError(errX, errY)
			} else {
				return 0, errors.New("division by zero")
			}
		}
	case *ast.UnaryExpr:
		if n.Op == token.SUB {
			val, err := Eval(n.X)
			return -val, err
		}
	case *ast.BasicLit:
		if n.Kind == token.INT || n.Kind == token.FLOAT {
			val, err := strconv.Atoi(n.Value)
			return float64(val), err
		}
	case *ast.ParenExpr:
		return Eval(n.X)
	}
	return 0, errors.New("unhandled node type")
}

// Returns the error if either err1 or err2 are non-nil.
// Returns err1 if both are non-nil.
func orError(err1 error, err2 error) error {
	if err1 != nil {
		return err1
	}
	return err2
}

func handleCalc(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	// Decode the request body into a go struct
	var request CalcRequest
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	// Generate an AST from the expression
	tree, err := parser.ParseExpr(request.Expression)
	if err != nil {
		fmt.Println("Bad request, handled error:", err)
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	// Evaluate the value from the AST
	result, err := Eval(tree)
	if err != nil {
		fmt.Println("Error during evaluation, probably caused by a bad request", err)
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	// Respond with the result
	response := CalcResponse{result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Write result to history file
	file_handling.WriteLine(fmt.Sprintf("%s;%f", request.Expression, result), historiesFileName)
}

func handleHistory(w http.ResponseWriter, numLines int) {
	lines, err := file_handling.GetLastLines(numLines, historiesFileName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	var calculations []Calculation
	for i := 0; i < len(lines); i++ {
		lineSplit := strings.Split(lines[i], ";")
		expr := lineSplit[0]
		res, err := strconv.ParseFloat(lineSplit[1], 64)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		calculations = append(calculations, Calculation{expr, res})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calculations)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(strings.TrimPrefix(r.URL.Path, apiRoot), "/")

	// Very simple router
	switch urlParts[0] {
	case "calc":
		if r.Method != "POST" {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if len(urlParts) > 1 && len(urlParts[1]) != 0 {
			http.Error(w, "404 not found", http.StatusNotFound) // Disallow /api/calc/*
			return
		}

		handleCalc(w, r)
		break
	case "history":
		if r.Method != "GET" {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if len(urlParts) > 2 && len(urlParts[2]) != 0 {
			http.Error(w, "404 not found", http.StatusNotFound) // Disallow /api/history/*/*
			return
		}

		numLines := defaultNumHistoryLines
		if len(urlParts) > 1 && len(urlParts[1]) > 0 {
			var err error
			numLines, err = strconv.Atoi(urlParts[1])
			if err != nil || numLines <= 0 {
				http.Error(w, "400 bad request", http.StatusBadRequest)
				return
			}
		}

		handleHistory(w, numLines)
		break
	case "coffee":
		if len(urlParts) > 1 && len(urlParts[1]) != 0 {
			http.Error(w, "404 not found", http.StatusNotFound) // Disallow /api/coffee/*
			return
		}
		http.Error(w, "418 I'm a teapot", http.StatusTeapot)
		break
	default:
		http.Error(w, "404 not found", http.StatusNotFound)
	}
}

func main() {
	// Create db file if it doesn't exist
	if _, err := os.Stat(historiesFileName); os.IsNotExist(err) {
		file, err := os.Create(historiesFileName)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	http.HandleFunc(apiRoot, apiHandler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
