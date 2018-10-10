package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const apiRoot = "/api/"

type CalcRequest struct {
	Expression string `json:"expression"`
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
	fmt.Println(request.Expression)
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
