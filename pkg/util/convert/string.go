package convert

import "strings"

// FormatJsonString format json string by remove slash "\""
func FormatJsonString(s string) string {
	return strings.ReplaceAll(s, "\"", "")
}
