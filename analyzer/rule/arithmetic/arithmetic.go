/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package arithmetic

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/object"
	"fluent/analyzer/queue"
	"fluent/ast"
	"fluent/filecode/types"
)

// AnalyzeArithmetic analyzes an arithmetic expression in the AST.
// It checks if the type of the expression can be inferred and if it matches the expected type.
// If the type cannot be inferred or does not match, it returns an error.
// Otherwise, it processes the expression and updates the expression queue.
//
// Parameters:
// - input: The AST node representing the arithmetic expression.
// - currentElement: The current expected pair being analyzed.
// - exprQueue: The queue of expected pairs to be processed.
//
// Returns:
// - An error3.Error indicating the result of the analysis.
func AnalyzeArithmetic(
	input *ast.AST,
	currentElement *queue.ExpectedPair,
	exprQueue *[]queue.ExpectedPair,
) error3.Error {
	expected := currentElement.Expected

	// Check if we can infer the type of the expression
	if expected.BaseType == "" {
		return error3.Error{
			Code:   error3.CannotInferType,
			Line:   input.Line,
			Column: input.Column,
		}
	}

	// Check if the expected is either a num or a dec
	if expected.BaseType != "num" && expected.BaseType != "dec" && expected.BaseType != "(Infer)" {
		return error3.Error{
			Code:       error3.TypeMismatch,
			Line:       input.Line,
			Column:     input.Column,
			Additional: []string{"num or dec", expected.BaseType},
		}
	}

	// Handle inferred types
	var candidateType types.TypeWrapper
	if expected.BaseType == "(Expected)" {
		candidateType = types.TypeWrapper{
			Children: &[]*types.TypeWrapper{},
		}
	} else {
		// Clone the expected element to avoid memory issues
		candidateType = types.TypeWrapper{
			PointerCount: expected.PointerCount,
			ArrayCount:   expected.ArrayCount,
			Children:     expected.Children,
			BaseType:     expected.BaseType,
			Trace:        expected.Trace,
			IsPrimitive:  expected.IsPrimitive,
		}
	}

	// Get the expression's children
	children := *input.Children

	// Push the candidate to determine the expression's type
	candidate := children[0]
	candidateElement := object.Object{
		Type: types.TypeWrapper{
			Children: &[]*types.TypeWrapper{},
		},
	}

	// Prevent nesting problems
	startAt := 0
	if candidate.Rule != ast.Expression {
		startAt = 1
	}

	for candidate.Rule == ast.Expression {
		newCandidate := (*candidate.Children)[0]

		if newCandidate.Rule == ast.Expression {
			candidate = newCandidate
			continue
		}

		if newCandidate.Rule != ast.ArithmeticExpression {
			break
		}

		candidate = newCandidate
	}

	*exprQueue = append(*exprQueue, queue.ExpectedPair{
		Expected:     &candidateType,
		Got:          &candidateElement,
		Tree:         candidate,
		IsArithmetic: true,
	})

	// Push the rest of the expression
	for i := startAt; i < len(children); i++ {
		element := children[i]

		// Skip sings
		if element.Rule == ast.ArithmeticSign {
			continue
		}

		*exprQueue = append(*exprQueue, queue.ExpectedPair{
			Expected: &candidateElement.Type,
			Got: &object.Object{
				Type: types.TypeWrapper{
					Children: &[]*types.TypeWrapper{},
				},
			},
			Tree: element,
		})
	}

	return error3.Error{}
}
