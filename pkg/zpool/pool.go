// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package zpool

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Property struct {
	Index int
	Value string
}

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
	for _, line := range lines {
		container := parseContainer(line)

		if container.Level != -1 {
			pool.Containers = append(pool.Containers, &container)
		}
	}

	return pool
}

func ListZpools() []string {
	cmd := append(cmdListPools, "name")
	return strings.Split(getOutput(cmd), "\n")
}

func GetProperties(name string, which string, props string) map[string]*Property {
	pulled := make(map[string]*Property)

	props = Sanitize(props)
	rawProps := strings.Split(props, ",")

	cmd := []string{ }
	if (which == "pool") {
		cmd = append(cmd, cmdListPools...)

	} else if (which == "dataset") {
		cmd = append(cmd, cmdListDatasets...)

	} else {
		log.Fatalf("Unknown item '%s' for GetProperties", which)
	}

	cmd = append(cmd, props, name)		// Append properties and pool name

	// This command returns a single line of output with properties delimited by tabs.
	// The order is determined by the properties passed to the -o flag.
	output := getOutput(cmd)
	parsed := strings.Split(output, "\t")

	for index, prop := range parsed {
		name := rawProps[index]
		cleaned := strings.Replace(prop, "\n", "", 1)		// The last property will have a newline at the end but to be safe, we'll clean every returned value

		pulled[name] = &Property {
			Index: index,
			Value: cleaned,
		}
	}

	return pulled
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

		pools = append(pools, ParsePool(name))
	}

	return pools
}

func ParsePool(name string) *Pool {
	name = Sanitize(name)
	cmd := append(cmdPoolStatus, name)
	out := getOutput(cmd)

	pool := ParseZpoolStatus(out)
	pool.Properties = GetProperties(name, "pool", "size,alloc,free,checkpoint,fragmentation,capacity,health,failmode,ashift")

	return pool
}

func GetVersion() string {
	out := getOutput(cmdGetVersion)

	// The first line is the zfs version, the second is the kernel module version
	version := strings.Split(out, "\n")[0]

	return version
}

func Encode(raw interface{}) []byte {
	encoded, err := json.Marshal(raw)

	if err != nil {
		log.Fatalf("Unable to marshal %#v to JSON: %s", raw, err)
	}

	return encoded
}

func parseContainer(line string) Container {
	info := strings.Fields(line)
	name := ""
	state := ""
	read := ""
	write := ""
	cksum := ""
	status := ""
	level := 0

	if len(info) == 0 {
		return Container {
			Level: -1,
		}
	}

	name = info[0]

	if len(info) > 1 {
		// cache entries only have one field, all others have (at least) 5
		state = info[1]
		read = info[2]
		write = info[3]
		cksum = info[4]
		level = countSpaces(line) / 2
	}

	if len(info) > 5 {
		// additional status information is available for this vdev member
		for i := 5; i < len(info); i++ {
			status += info[i] + " "
		}

		status = strings.TrimSpace(status)
	}

	return Container {
		Name:   name,
		State:  state,
		Read:   read,
		Write:  write,
		Cksum:  cksum,
		Status: status,
		Level:  level,
	}
}
