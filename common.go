package main

import (
	"strings"
	"os"
)

/*
	Attempt to determine severity of message from keywords
*/
func calculateSeverity(message string) int64 {
	var severity int64 = 1
	lower := strings.ToLower(message)

	// Check warning keywords
	for _, keyword := range global_cfg.Client.WarningKeywords {
		if (strings.Contains(lower, keyword)) {
			severity = 2
			break
		}
	}

	// Check error keywords
	for _, keyword := range global_cfg.Client.ErrorKeywords {
		if (strings.Contains(lower, keyword)) {
			severity = 3
			break
		}
	}

	return severity
}

/*
	Check if file exists
*/
func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if (os.IsNotExist(err)) {
        return false
    }
    return !info.IsDir()
}