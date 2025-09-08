// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>
package main

import (
	"fmt"

	"github.com/ericlagergren/decimal"
	"github.com/fatih/color"
)

type (
	// stackType holds the representation of the RPN stack. It contains
	// two stacks, "list" (the main stack), and "savedList", which is
	// used to save the stack and later restore it in case of error.
	stackType struct {
		list      []*decimal.Big
		savedList []*decimal.Big
	}
)

// save saves the current stack in a separate structure.
func (x *stackType) save() {
	x.savedList = append([]*decimal.Big{}, x.list...)
}

// restore restores the saved stack back into the main one.
func (x *stackType) restore() {
	x.list = append([]*decimal.Big{}, x.savedList...)
}

// push adds a new element to the stack.
func (x *stackType) push(n ...*decimal.Big) {
	x.list = append(x.list, n...)
}

// clear clears the stack.
func (x *stackType) clear() {
	x.list = []*decimal.Big{}
}

// top returns the topmost element on the stack (without popping it).
func (x *stackType) top() *decimal.Big {
	if len(x.list) == 0 {
		return big()
	}
	return x.list[len(x.list)-1]
}

// printTop displays the top of the stack using the base indicated.
func (x *stackType) printTop(ctx decimal.Context, base, decimals int) {
	color.Cyan("= %s", formatNumber(ctx, x.top(), base, decimals, false))
}

// print displays the contents of the stack using the base indicated.
func (x *stackType) print(ctx decimal.Context, base, decimals int) {
	last := len(x.list) - 1

	fmt.Println(bold("===== Stack ====="))
	for ix := last; ix >= 0; ix-- {
		tag := fmt.Sprintf("%2d", ix)
		switch ix {
		case last:
			tag = " x"
		case last - 1:
			tag = " y"
		}
		fmt.Printf("%s: %s\n", tag, formatNumber(ctx, x.list[ix], base, decimals, false))
	}
}
