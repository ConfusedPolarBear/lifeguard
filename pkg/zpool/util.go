// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package zpool

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// Sanitizes a string so it is safe to use as a shell argument. Only characters that are valid in zfs datasets, snapshots or properties are permitted.
func Sanitize(raw string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\-_:\.%,]`)
	return re.ReplaceAllString(raw, "")
}

func getOutput(raw []string) string {
	cmd := exec.Command(raw[0], raw[1:]...)
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatalf("Unable to exec %s: %s", raw, err)
	}

	return string(stdout)
}

// Utility command for grabbing contents to the right of defined string
func grabRightOf(str string, find string) string {
	return str[len(find):]
}

// Utility function to find a string in the first 7 characters of the haystack as the headers of the zpool status command are padded to a constant display width
// Example: "state: ONLINE"
func hasHeader(haystack string, needle string) bool {
	first := haystack[0:7]
	return strings.Index(first, needle) != -1
}

// Utility function to turn "examplekey: Example value with spaces" into "examplekey" and "Example value with spaces"
func toMap(raw string) (string, string) {
	key := strings.Split(raw[0:7], ":")[0]
	length := len(key + ": ")
	value := raw[length:]

	return key, value
}

// Utility function to use the provided value or the provided default if value is empty
func useDefault(value string, def string) string {
	if value == "" {
		return def
	}

	return value
}

// Utility function to count the number of leading spaces in str.
func countSpaces(str string) int {
	count := 0
	for  _, current := range str {
		if current == ' ' {
			count++
		} else if current == 9 {
			// ignore tabs
			continue
		} else {
			break
		}
	}

	return count
}
