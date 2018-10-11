package main

import (
	"fmt"
	"go/parser"
	"testing"
)

func TestEval(t *testing.T) {
	// Good expressions
	expressions := []string{
		"-1 * (2 * 6 / 3)",
		"2 + 2",
		"2        +         (2)",
		"1 * 1 * 1 * 1 * 1 * 0",
		"1 * (2 * (3 * (4 * (0))))",
		"-(-(-(-1)))",
		"(1 + 2) * 3",
		"1 + 2 * 3",
		"3 ///// 1", // I have decided to allow expressions like this
	}
	expectations := []float64{-4, 4, 4, 0, 0, 1, 9, 7, 3}
	for i := 0; i < len(expressions); i++ {
		tree, err := parser.ParseExpr(expressions[i])
		if err != nil {
			fmt.Println("Unexpected error when parsing", expressions[i], ":", err)
			t.Fail()
			continue
		}

		result, err := Eval(tree)
		if err != nil {
			fmt.Println("Unexpected error when evaluating", expressions[i], ":", err)
			t.Fail()
			continue
		}

		fmt.Printf("%s -> %f", expressions[i], result)
		if result == expectations[i] {
			fmt.Printf(" (OK)\n")
		} else {
			fmt.Printf(" (FAIL) (expected %f)\n", expectations[i])
			t.Fail()
		}
	}

	// Bad expressions
	expressions = []string{
		"2 ++ 2",
		"2 * (5 + 1",
		"3 / )(5)",
		"7 / 0",
	}
	for i := 0; i < len(expressions); i++ {
		tree, err := parser.ParseExpr(expressions[i])
		if err != nil {
			fmt.Println("Got expected error when parsing", expressions[i], ":", err, "(OK)")
			continue
		}

		result, err := Eval(tree)
		if err != nil {
			fmt.Println("Got expected error when evaluating", expressions[i], ":", err, "(OK)")
			continue
		}

		fmt.Println("Bad expression", expressions[i], "evaluated to", result, "without returning error (FAIL)")
		t.Fail()
	}
}
