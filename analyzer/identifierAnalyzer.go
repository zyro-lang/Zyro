package analyzer

import (
	"surf/ast"
	"surf/code"
	"surf/core/engine/args"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
	"surf/token"
	"surf/tokenUtil"
)

// AnalyzeIdentifier analyzes the given identifier
// and next tokens
func AnalyzeIdentifier(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]*code.SurfMod,
	startAt *int,
	lastValue *object.SurfObject,
	isArithmetic *bool,
	isFunCall *bool,
) {
	firstToken := statement[0]

	// Check if the identifier is a variable
	// or a function call
	variable, varFound := variables.Load(firstToken.GetValue())
	function, funFound, sameFile := ast.LocateFunction(*functions, firstToken.GetFile(), firstToken.GetValue())

	if !varFound && !funFound {
		logger.TokenError(
			firstToken,
			"Undefined reference to variable "+firstToken.GetValue(),
			"The variable "+firstToken.GetValue()+" was not found in the current scope",
			"Declare the variable in the current scope",
		)
	}

	if funFound {
		if !sameFile && !function.IsPublic() {
			logger.TokenError(
				firstToken,
				"Function "+firstToken.GetValue()+" is not public",
				"Move the function to the current file or make it public",
			)
		}

		statementLen := len(statement)
		if statementLen < 3 || statement[1].GetType() != token.OpenParen {
			logger.TokenError(
				firstToken,
				"Invalid function call",
				"A function call must be followed by parentheses",
				"Add parentheses to the function call",
			)
		}

		// Parse and check the arguments
		// Extract all the tokens of the function invocation
		call := tokenUtil.ExtractTokensBefore(
			statement,
			token.CloseParen,
			true,
			token.OpenParen,
			token.CloseParen,
			true,
		)

		// ExtractTokensBefore also checks the end closing parenthesis
		// is also met, so no need to check it here

		argumentsRaw, skipped := args.SplitArgs(call)
		arguments := make([]object.SurfObject, len(argumentsRaw))

		for i, argument := range argumentsRaw {
			argValue := AnalyzeStatement(
				argument,
				variables,
				functions,
				mods,
			)

			arguments[i] = argValue
		}

		*startAt += skipped

		*lastValue = AnalyzeFun(
			function,
			functions,
			mods,
			firstToken,
			true,
			stack.NewStack(),
			arguments...,
		)
	} else {
		*lastValue = variable
		*startAt += 1
	}

	if (*lastValue).GetType() == object.IntType || (*lastValue).GetType() == object.DecimalType {
		*isArithmetic = true
	}

	*isFunCall = funFound
}
