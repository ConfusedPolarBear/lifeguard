// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package zpool

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
)

func ExecWithInput(raw []string, stdin []byte) (string, string, error) {
	return execInternal(raw, stdin)
}

func Exec(raw []string) (string, string, error) {
	return execInternal(raw, []byte(""))
}

func MustExec(raw []string) string {
	stdout, _, err := Exec(raw)

	if err != nil {
		log.Fatalf("Unable to exec %s: %s", raw, err)
	}

	return stdout
}

func execInternal(raw []string, stdin []byte) (string, string, error) {
	var stdout, stderr bytes.Buffer

	raw = append([]string { config.GetString("exec.timeout_path"), config.GetString("exec.timeout") }, raw...)
	cmd := exec.Command(raw[0], raw[1:]...)

	if config.GetBool("debug.exec") {
		log.Printf("Executing command: %v", cmd)
	}

	cmd.Stdin = bytes.NewBuffer(stdin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", string(stderr.Bytes()), err
	}

	return string(stdout.Bytes()), string(stderr.Bytes()), nil
}

// Sanitizes a string so it is safe to use as a shell argument. Only characters that are valid in zfs datasets, snapshots or properties are permitted.
func Sanitize(raw string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\-_:\.%,]`)
	return re.ReplaceAllString(raw, "")
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
