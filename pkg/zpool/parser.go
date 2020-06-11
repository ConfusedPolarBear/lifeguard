// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package zpool

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var cmdListPools  = []string { "/sbin/zpool", "list", "-H", "-o", "name" }
var cmdPoolStatus = []string { "/sbin/zpool", "status" }

func ParseZpoolStatus(raw string) *Pool {
	var pool *Pool

	inConfig := false
	config := ""

	var processed []string
	poolMap := make(map[string]string)

	lines := strings.Split(raw, "\n")

	// Loop through stdout and combine wrapped lines (lines where the current line is a continuation of the previous one)
	for _, line := range lines {
		if len(line) <= 1 {
			continue
		}

		// Wrapped lines do not have a colon at the start
		if hasHeader(line, ":") {
			line = strings.TrimSpace(line)

			if hasHeader(line, "config:") {
				inConfig = true
				continue

			} else if hasHeader(line, "errors:") && inConfig {
				// The errors section always appears immediately after the config
				inConfig = false
			}

			processed = append(processed, line)

		} else {
			if inConfig {
				config += line + "\n"

			} else {
				line = strings.TrimSpace(line)
				processed[len(processed) - 1] += "\n" + line
			}
		}
	}

	for _, l := range processed {
		key, value := toMap(l)
		poolMap[key] = value
	}
	poolMap["config"] = config
	poolMap["raw"] = raw

	pool = &Pool {
		Name:   poolMap["pool"],
		State:  poolMap["state"],
		Status: useDefault(poolMap["status"], "OK"),
		Action: useDefault(poolMap["action"], "No action needed"),
		See:    poolMap["see"],
		Scan:   poolMap["scan"],
		Errors: poolMap["errors"],
		Raw:    poolMap["raw"],
	}

	if os.Getenv("LIFEGUARD_DEBUG") == "1" {
		log.Printf("==================== Processed zpool output =======================\n")
		for key, value := range poolMap {
			log.Printf("'%s' => '%s'\n", key, value)
		}
	}

	lines = strings.Split(poolMap["config"], "\n")

	// The first line is the header and should be skipped
	lines[0] = ""
	initialLevel := -1

	for _, line := range lines {
		// TODO: replace with map?
		info := strings.Fields(line)
		name := ""
		state := ""
		read := ""
		write := ""
		cksum := ""
		status := ""
		level := 0

		if len(info) == 0 {
			continue
		}

		name = info[0]
		if initialLevel == -1 {
			initialLevel = countSpaces(line)
		}

		if len(info) > 1 {
			// cache entries only have one field, all others have (at least) 5
			state = info[1]
			read = info[2]
			write = info[3]
			cksum = info[4]
			level = (countSpaces(line) - initialLevel) / 2
		}

		if len(info) > 5 {
			// additional status information is available for this vdev member
			for i := 5; i < len(info); i++ {
				status += info[i] + " "
			}

			status = strings.TrimSpace(status)
		}

		pool.Containers = append(pool.Containers, &Container {
			Name:   name,
			State:  state,
			Read:   read,
			Write:  write,
			Cksum:  cksum,
			Status: status,
			Level:  level,
		})
	}

	return pool
}

func ListZpools() []string {
	return strings.Split(getOutput(cmdListPools), "\n")
}

func ParseAllPools() []*Pool {
	var pools []*Pool
	names := ListZpools()

	if len(names) > 0 && names[0] == "no pools available" {
		return pools
	}

	for _, name := range names {
		if name == "" {
			continue
		}

		name = Sanitize(name)

		cmd := append(cmdPoolStatus, name)
		out := getOutput(cmd)

		pools = append(pools, ParseZpoolStatus(out))
	}

	return pools
}

func GetPools() []byte {
	pools := ParseAllPools()
	encoded, err := json.Marshal(pools)

	if err != nil {
		log.Fatalf("Unable to marshal pool to JSON: %s", err)
	}

	return encoded
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

// Sanitizes a string so it is safe to use as a shell argument. Only letters, numbers, dashes and underscore are allowed
func Sanitize(raw string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\-_]`)
	return re.ReplaceAllString(raw, "")
}

// Utility function to count the number of leading spaces in str.
func countSpaces(str string) int {
	count := 0
	for  _, current := range str {
		if current == ' ' {
			count++
		} else {
			break
		}
	}

	return count
}
