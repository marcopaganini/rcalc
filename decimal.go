// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>
package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ericlagergren/decimal"
)

// big returns a new *decimal.Big
func big() *decimal.Big {
	return decimal.WithContext(decimal.Context128)
}

// bigUint returns a new *decimal.Big from an uint64.
func bigUint(n uint64) *decimal.Big {
	return big().SetUint64(n)
}

// bigFloat returns a new *decimal.Big from a string
// Using a float64 here will introduce rounding errors.
func bigFloat(s string) *decimal.Big {
	r, _ := big().SetString(s)
	return r
}

// commafWithDigits idea comes from the humanize library, but was modified to
// work with decimal numbers.
func commafWithDigits(n *decimal.Big, decimals int) string {
	// Make a copy so we won't modify the original value (passed by pointer).
	v := big().Copy(n)

	buf := &bytes.Buffer{}
	if v.Sign() < 0 {
		buf.Write([]byte{'-'})
		// Make v positive
		v.SetSignbit(false)
	}

	comma := []byte{','}

	f := fmt.Sprintf("%%.%df", decimals)
	parts := strings.Split(fmt.Sprintf(f, v), ".")

	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return stripTrailingDigits(buf.String(), decimals)
}

func stripTrailingDigits(s string, digits int) string {
	// Remove insignificant zeroes after period (if any).
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}

	// Trim string to .<digits>
	if i := strings.Index(s, "."); i >= 0 {
		if digits <= 0 {
			return s[:i]
		}
		i++
		if i+digits >= len(s) {
			return s
		}
		return s[:i+digits]
	}
	return s
}

// formatNumber formats the number using base and decimals. For bases different
// than 10, non-integer floating numbers are truncated.
func formatNumber(ctx decimal.Context, n *decimal.Big, base, decimals int, single bool) string {
	// Print NaN without suffix numbers.
	if n.IsNaN(0) {
		return strings.TrimRight(fmt.Sprint(n), "0123456789")
	}
	if n.IsInf(0) {
		return fmt.Sprint(n)
	}

	// clean = double as ascii, without non-significant decimal zeroes.
	f := fmt.Sprintf("%%.%df", decimals)
	clean := stripTrailingDigits(fmt.Sprintf(f, n), decimals)

	var (
		n64    uint64
		suffix string
	)

	buf := &bytes.Buffer{}
	if base != 10 {
		// For negative numbers, prefix them with a minus sign and
		// force them to be positive.
		if n.Signbit() {
			buf.Write([]byte{'-'})
			n.SetSignbit(false)
		}
		// Truncate floating point numbers to their integer representation.
		if !n.IsInt() {
			suffix = fmt.Sprintf(" (truncated from %s)", clean)
			ctx.Floor(n, n)
		}
		// Non-base 10 uses uint64s.
		var ok bool
		n64, ok = n.Uint64()
		if !ok {
			return "Invalid number: non decimal base only supports uint64 numbers."
		}
	}

	switch base {
	case 2:
		fmt.Fprintf(buf, "0b%b%s", n64, suffix)
	case 8:
		fmt.Fprintf(buf, "0%o%s", n64, suffix)
	case 16:
		fmt.Fprintf(buf, "0x%x%s", n64, suffix)
	default:
		h := commafWithDigits(n, decimals)
		// Only print humanized format when it differs from original value, and not in single-command mode
		if h != clean && !single {
			suffix = " (" + h + ")"
		}
		buf.WriteString(clean + suffix)
	}

	return buf.String()
}
