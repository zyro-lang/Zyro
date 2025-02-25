/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package cli

import (
	"fluent/ansi"
	"fluent/filecode/converter"
	"fluent/ir"
	"fluent/logger"
	"fluent/state"
	"fluent/util"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"
)

var includePath = fmt.Sprintf("%s/include", converter.StdPath)
var includePathPOSIX = fmt.Sprintf("%s/posix", includePath)
var includePathWin = fmt.Sprintf("%s/win", includePath)
var fluentExtensionRegex = regexp.MustCompile("\\.fluent$")
var isWindows = runtime.GOOS == "windows"
var isPOSIX = !isWindows

// BuildCommand compiles the given Fluent project into an executable
func BuildCommand(context *cli.Command) {
	fmt.Println(ansi.Colorize(ansi.BoldBrightYellow, "⚠️ Checking if fluentc is installed...."))

	// Invoke a system command to check if fluentc is installed
	cmd := exec.Command("fluentc", "--help")
	err := cmd.Run()

	if err != nil {
		logger.Error("The Fluent Compiler is not installed.")
		logger.Info(
			"Please install it by downloading the necessary",
			"binaries from the official repository.",
		)
		os.Exit(1)
	}

	fileCodes, fileCodesMap := CheckCommand(context)
	// Retrieve the path from the context
	userPath := util.GetDir(context.Args().First())

	// Use a global builder to build the whole program into a single IR file
	globalBuilder := strings.Builder{}

	// Keep a counter of all the file codes that have been processed
	fileCodeCount := 0

	// A map of already-defined values used for tracing lines and columns
	traceCounters := make(map[int]int)
	// Keep track of used strings (Saved in reserved spaces of memory)
	usedStrings := make(map[string]string)
	// Used to store precomputed counters for functions' and
	// modules' names
	nameCounters := make(map[string]map[string]string)
	globalCounter := 0
	// Save in a map the files that have an external
	// implementation to avoid recompiling them
	externalImpl := make(map[string]bool)

	// Precompute the counters for the names
	for _, fileCode := range fileCodes {
		// Check if this file has an external implementation
		if strings.HasPrefix(fileCode.Path, converter.StdPath) {
			var relativePath string

			// Check for POSIX-Compliant systems
			if isPOSIX {
				relativePath = strings.Replace(
					fileCode.Path,
					converter.StdPath,
					includePathPOSIX,
					1,
				)
			} else {
				relativePath = strings.Replace(
					fileCode.Path,
					converter.StdPath,
					includePathWin,
					1,
				)
			}

			relativePath = fluentExtensionRegex.ReplaceAllString(relativePath, ".ll")
			if util.FileExists(relativePath) {
				fileName := util.FileName(&fileCode.Path)
				externalImpl[fileCode.Path] = true

				fmt.Println(
					ansi.Colorize(
						ansi.BoldBrightYellow,
						fmt.Sprintf(
							"🔂 Skipped %s (External impl available)",
							fileName,
						),
					),
				)

				// Add the std instruction to the global builder
				globalBuilder.WriteString("link ")
				globalBuilder.WriteString(relativePath)
				globalBuilder.WriteString("\n")
				continue
			}
		}

		// Make sure the map is initialized
		nameCounters[fileCode.Path] = make(map[string]string)
		// Get the stored map
		nameCounter := nameCounters[fileCode.Path]

		for _, fun := range fileCode.Functions {
			nameCounter[fun.Name] = fmt.Sprintf("x%d", globalCounter)
			globalCounter++
		}

		for _, mod := range fileCode.Modules {
			nameCounter[mod.Name] = fmt.Sprintf("x%d", globalCounter)
			globalCounter++
		}
	}

	for _, fileCode := range fileCodes {
		// Skip the file if it has an external implementation
		if externalImpl[fileCode.Path] {
			continue
		}

		fileCodeCount++
		fileName := util.FileName(&fileCode.Path)

		// Emit a building state
		state.Emit(state.Building, fileName)

		fileIr := ir.BuildIr(
			fileCode,
			fileCodesMap,
			fileCodeCount,
			&traceCounters,
			&usedStrings,
			// Prevent copying the map every time
			// by passing a reference to the map
			&nameCounters,
			nameCounters[fileCode.Path],
		)
		// Write the IR to the global builder
		globalBuilder.WriteString(fileIr)
		globalBuilder.WriteString("\n")
		state.PassAllSpinners()
	}

	// Get the pwd
	pwd, err := os.Getwd()

	if err != nil {
		logger.Error("Could not get the current working directory.")
		os.Exit(1)
	}

	// Get the out directory path
	outDir := path.Join(pwd, userPath, "out")

	// Make sure the output directory exists
	if !util.DirExists(outDir) {
		err := os.Mkdir(outDir, os.ModePerm)

		if err != nil {
			logger.Error("Could not create the output directory.")
			os.Exit(1)
		}
	}

	// Write the global IR to a file
	globalIrPath := path.Join(outDir, "program.flc")
	outPath := path.Join(outDir, "out")

	// Add an .exe extension if the user is on Windows
	if !isPOSIX {
		outPath += ".exe"
	}

	err = os.WriteFile(globalIrPath, []byte(globalBuilder.String()), os.ModePerm)

	if err != nil {
		logger.Error("Could not write the Fluent IR to a file.")
		os.Exit(1)
	}

	fmt.Println(ansi.Colorize(ansi.BoldBrightYellow, "⚠️ Invoking fluent backend...."))
	fmt.Println(ansi.Colorize(ansi.BrightBlack, "⚠️ The output you will see from now on is coming from the fluentc command."))

	// Invoke the fluentc backend
	cmd = exec.Command("fluentc", "-o", outPath, globalIrPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Errors are handled by the compiler backend
	_ = cmd.Run()
}
