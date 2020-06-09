// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package zpool

import (
	// "encoding/json"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var cmdListPools  = []string { "/sbin/zpool", "list", "-H", "-o", "name" }
var cmdPoolStatus = []string { "/sbin/zpool", "status" }

func ParseZpoolStatus(raw string) *PoolStatus {
	var pool *PoolStatus

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

	pool = &PoolStatus{Name: poolMap["pool"]}
	pool.State = poolMap["state"]
	pool.Status = useDefault(poolMap["status"], "OK")
	pool.Action = useDefault(poolMap["action"], "No action needed")
	pool.See = poolMap["see"]
	pool.Scan = poolMap["scan"]
	pool.Errors = poolMap["errors"]
	pool.Raw = poolMap["raw"]

	if os.Getenv("DEBUG") == "1" {
		log.Printf("==================== Processed zpool output =======================\n")
		for key, value := range poolMap {
			log.Printf("'%s' => '%s'\n", key, value)
		}
	}

	lines = strings.Split(poolMap["config"], "\n")

	// The first line is the header and should be skipped
	lines[0] = ""
	for _, line := range lines {
		if line == "" {
			continue
		}

		info := strings.Fields(line)
		status := "none"
		if len(info) > 5 {
			status = ""
			for i := 5; i < len(info); i++ {
				status += info[i] + " "
			}
		}

		pool.Containers = append(pool.Containers, &ContainerStatus{
			Name:   info[0],
			State:  info[1],
			Read:   info[2],
			Write:  info[3],
			Cksum:  info[4],
			Status: status})
	}

	return pool
}

func ListZpools() []string {
	return strings.Split(getOutput(cmdListPools), "\n")
}

func ParseAllPools() []*PoolStatus {
	var pools []*PoolStatus
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

