package main

import (
	"strings"
)

func calculateSeverity(message string) int64 {
	var severity int64 = 1
	lower := strings.ToLower(message)

	if (strings.Contains(lower, "warn") || strings.Contains(lower, "warning")) {
		severity = 2
	}
	if (strings.Contains(lower, "fatal") || strings.Contains(lower, "error")) {
		severity = 3
	}

	return severity
}