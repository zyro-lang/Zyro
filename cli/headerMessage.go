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
	"fmt"
)

// ShowHeaderMessage shows the header message of the Fluent CLI
func ShowHeaderMessage() {
	fmt.Println(ansi.Colorize(ansi.BoldBrightBlue, "The Fluent Programming Language"))
	fmt.Println(ansi.Colorize(ansi.BrightBlack, "A blazingly fast programming language"))
	fmt.Println()
}
