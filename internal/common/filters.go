// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package common

import "strings"

// QuoteFilterValue wraps s in double quotes for an ISC filter expression,
// escaping backslashes and embedded quotes.
func QuoteFilterValue(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}
