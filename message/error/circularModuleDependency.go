/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package error

import (
	"fluent/logger"
	"strings"
)

func CircularModuleDependency(chain string) string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("This module depends on its own"))
	builder.WriteString(
		logger.BuildHelp(
			"This module has a property that is either the same module or",
			"another module that depends on this one.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0020",
			"Full dependency chain:",
		),
	)

	builder.WriteString(chain)
	builder.WriteString(
		logger.BuildInfo(
			"Full details:",
		),
	)

	return builder.String()
}
