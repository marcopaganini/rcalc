// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/ericlagergren/decimal"
	"github.com/fatih/color"
)

var (
	// Build is filled by go build -ldflags during build.
	Build        string
	programTitle = "rpn - a simple CLI RPN calculator"

	// These are functions to be used to print in color.
	errorMsg = color.New(color.FgRed).SprintFunc()
	warnMsg  = color.New(color.FgMagenta).SprintFunc()
	bold     = color.New(color.Bold).SprintFunc()
)

// atof takes a string as an argument and return a decimal object representing
// that string. Strings starting in 0x or 0X are treated as hex strings.
// Strings starting in o or 0 are treated as octal strings. Non decimal strings
// are converted to a uint64 intermediate representation and thus limited to
// how much a uint64 can hold.
func atof(s string) (*decimal.Big, error) {
	base := 10
	switch {
	case (strings.HasPrefix(s, "0b") || strings.HasPrefix(s, "0B")) && len(s) > 2:
		s = s[2:]
		base = 2
	case (strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X")) && len(s) > 2:
		s = s[2:]
		base = 16
	// Numbers starting with 0 must account for 0.xx fractional numbers not
	// being octal numbers.
	case (strings.HasPrefix(s, "0") || strings.HasPrefix(s, "o")) && !strings.HasPrefix(s, "0.") && len(s) > 1:
		s = s[1:]
		base = 8
	}

	if base == 10 {
		var d decimal.Big
		if _, ok := d.SetString(s); !ok || d.IsNaN(0) {
			return nil, errors.New("unable to convert number")
		}
		return &d, nil
	}

	// Non-base 10 numbers are limited to uint64 sizes.
	ret, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		return nil, err
	}
	return bigUint(ret), nil
}

// calc contains the bulk of the calculator code. It takes a stack and an
// optional string argument. If string the string is not empty, it executes the
// oeprations in the string and returns. If the string is empty, it enters a
// readline loop accepting commands from the user.
func calc(stack *stackType, cmd string) error {
	// Wait for entry until Ctrl-D or q is issued
	var (
		line string
		err  error
		rl   *readline.Instance
	)

	ctx := decimal.Context128

	// Single command execution?
	single := (cmd != "")

	// Operations
	ops := newOpsType(ctx, stack)
	opmap := ops.opmap()

	if !single {
		rl, err = readline.New("> ")
		if err != nil {
			log.Fatal(err)
		}
		defer rl.Close()
	}

	// Remove all extraneous characters from the input. This will silently
	// remove undesirable formatting characters, making cut/paste operations
	// simpler. If you add a new operation as a single special character, make
	// sure it's represented here.
	cleanRe := regexp.MustCompile(`[^-+./*%^=[:alnum:]\s]`)

	for {
		// Save a copy of the stack so we can restore it to the previous state
		// before this line was processed (in case of errors.)
		stack.save()

		if ops.debug {
			stack.print(ctx, ops.base, ops.decimals)
		}

		// By default, use the passed command. If no command, initialize readline.
		line = cmd
		if !single {
			line, err = rl.Readline()
			if err != nil { // io.EOF
				break
			}
		}
		// Comment?
		if strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimSpace(line)
		line = cleanRe.ReplaceAllString(line, "")

		// Split into fields and process
		autoprint := false
		for _, token := range strings.Fields(line) {
			// Check operator map
			handler, ok := opmap[token]
			if ok {
				results, remove, err := operation(handler, stack)
				if err != nil {
					if single {
						return err
					}
					fmt.Printf(errorMsg("ERROR: %v\n"), err)
					stack.restore()
					break
				}
				// If the particular handler does not ignore results from the
				// function, set autoprint to true. This will cause the top of
				// the stack results to be printed.
				autoprint = (len(results) > 0 || remove > 0)

				if !single {
					// Set readline prompt based on base and degrees/radian mode.
					switch {
					case ops.degmode:
						rl.SetPrompt("deg> ")
					case ops.base == 10:
						rl.SetPrompt("> ")
					case ops.base == 8:
						rl.SetPrompt("oct> ")
					case ops.base == 16:
						rl.SetPrompt("hex> ")
					case ops.base == 2:
						rl.SetPrompt("bin> ")
					}
				}
				continue
			}

			// Help
			if token == "help" || token == "h" || token == "?" {
				if err := ops.help(); err != nil {
					fmt.Println(errorMsg(err))
				}
				continue
			}

			if token == "quit" || token == "exit" || token == "q" {
				fmt.Printf("Bye.\n")
				os.Exit(0)
			}

			// At this point, it's either a number or not recognized.
			// If anything fails, restore stack and stop token processing.
			n, err := atof(token)
			if err != nil {
				fmt.Printf(errorMsg("Not a number or operator: %q.\n"), token)
				fmt.Println(errorMsg("Use \"help\" for online help."))
				stack.restore()
				break
			}
			// Valid number
			stack.push(n)
			continue
		}

		if autoprint {
			if single {
				fmt.Println(formatNumber(ctx, stack.top(), ops.base, ops.decimals, true)) // plain print to stdout
			} else {
				stack.printTop(ctx, ops.base, ops.decimals) // pretty print to terminal
			}
		}

		// Break after the first iteration if a command is passed.
		if single {
			break
		}
	}
	return nil
}

func main() {
	stack := &stackType{}

	if err := calc(stack, strings.Join(os.Args[1:], " ")); err != nil {
		log.Fatal(err)
	}
}
