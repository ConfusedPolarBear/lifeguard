// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package zpool

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"regexp"

	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
	"github.com/ConfusedPolarBear/lifeguard/pkg/crypto"
	"github.com/ConfusedPolarBear/lifeguard/pkg/notifications"
	"github.com/ConfusedPolarBear/lifeguard/pkg/structs"
)

var IsTest = false
var poolHistory = make(map[string]*structs.Pool)

func ParseZpoolStatus(raw string) *structs.Pool {
	var pool *structs.Pool

	inConfig := false
	poolConfig := ""

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

			} else if hasHeader(line, "errors:") {
				// The errors section always appears immediately after the config
				inConfig = false
			}

			processed = append(processed, line)

		} else {
			if inConfig {
				poolConfig += line + "\n"

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
	poolMap["config"] = poolConfig
	poolMap["raw"] = raw

	pool = &structs.Pool {
		Name:    poolMap["pool"],
		State:   poolMap["state"],
		Status:  useDefault(poolMap["status"], "OK"),
		Action:  useDefault(poolMap["action"], "No action needed"),
		See:     poolMap["see"],
		Scan:    poolMap["scan"],
		Scanned: 0.0,
		Errors:  poolMap["errors"],
		Raw:     poolMap["raw"],
	}

	// Search for a percent in the scan output and if found, save the scan percentage to the scanned field
	// This regex searches for (any numbers or a period) followed by a percent sign.
	percentRegex := regexp.MustCompile("[0-9.]+%")
	if percentRegex.MatchString(pool.Scan) {
		raw := percentRegex.FindString(pool.Scan)
		raw = strings.ReplaceAll(raw, "%", "")

		pool.Scanned, _ = strconv.ParseFloat(raw, 64)

		// Determine if the pool scan is paused
		if strings.Index(pool.Scan, "paused") != -1 {
			pool.ScanPaused = true
		}
	}

	if IsTest || config.GetBool("debug.parse", false) {
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

	// Check for pool state changes and send notifications as needed
	name := pool.Name
	notifications.UpdatePoolState(name, pool, poolHistory[name])
	
	// Save the current state
	poolHistory[name] = pool

	return pool
}

func ListZpools() []string {
	cmd := append(cmdListPools, "name")
	return strings.Split(MustExec(cmd), "\n")
}

func GetProperties(name string, which string, filter string, props string) []map[string]*structs.Property {
	var pulled []map[string]*structs.Property

	props = Sanitize(props)
	rawProps := strings.Split(props, ",")

	if filter != "filesystem" && filter != "snapshot" && filter != "" {
		log.Fatalf("Unknown filter type: %s", filter)
	}

	cmd := []string{ }
	if (which == "zpool") {
		cmd = append(cmd, cmdListPools...)

	} else if (which == "zfs") {
		cmd = append(cmd, cmdListDatasets...)

	} else {
		log.Fatalf("Unknown command '%s' for GetProperties", which)
	}

	cmd = append(cmd, props, name)		// Append properties and pool name
	if (which == "zfs" && filter != "") {
		// only the zfs command supports filtering by type and recursion
		cmd = append(cmd, "-t", filter, "-r")
	}

	// This command returns a single line of output with properties delimited by tabs.
	// The order is determined by the properties passed to the -o flag.
	output := MustExec(cmd)

	if IsTest || config.GetBool("debug.parse", false) {
		log.Printf("Raw output of %v: '%s'", cmd, output)
	}

	lines := strings.Split(output, "\n")

	for number, line := range lines {
		if len(line) <= 1 {
			continue
		}

		parsed := strings.Split(line, "\t")

		for index, prop := range parsed {
			name := rawProps[index]
			cleaned := strings.Replace(prop, "\n", "", 1)		// The last property will have a newline at the end but to be safe, we'll clean every returned value

			if len(pulled) <= number {
				blank := make(map[string]*structs.Property)
				pulled = append(pulled, blank)
			}

			/* Since names can't be sanitized without potentially breaking them, an HMAC is calculated over them.
			 * When a zfs command needs to be run with this unsanitized input, the API only accepts the HMAC as
			 * the name. Before executing the command, the HMAC looked up in RAM to ensure it is a value we generated.
			 * This allows the API to work with unsanitized text in a safe manner.
			 */
			hmac := ""
			if name == "name" {
				hmac = crypto.GenerateHMAC(cleaned)
			}

			pulled[number][name] = &structs.Property {
				Name:  name,
				Value: cleaned,
				HMAC:  hmac,
			}
		}
	}

	return pulled
}

func ParseAllPools() []*structs.Pool {
	var pools []*structs.Pool
	names := ListZpools()

	if len(names) > 0 && names[0] == "no pools available" {
		return pools
	}

	for _, name := range names {
		if name == "" {
			continue
		}

		pools = append(pools, ParsePool(name, false))
	}

	return pools
}

func ParsePool(name string, includeChildren bool) *structs.Pool {
	name = Sanitize(name)
	cmd := append(cmdPoolStatus, name)
	out := MustExec(cmd)

	pool := ParseZpoolStatus(out)
	pool.Properties = GetProperties(name, "zpool", "", config.GetString("properties.pool", structs.DefaultProperties["pool"]))[0]

	/*
	 * This is optional since parsing all snapshots is expensive if many are present.
	 * In one test with 126 snapshots (created automatically with sanoid), the API response time went
	 * from less than 50 ms on average to 260 ms.
	 */
	if includeChildren {
		pool.Datasets   = GetProperties(name, "zfs", "filesystem", config.GetString("properties.dataset", structs.DefaultProperties["dataset"]))
		pool.Snapshots  = GetProperties(name, "zfs", "snapshot", config.GetString("properties.snapshot", structs.DefaultProperties["snapshot"]))
	}

	return pool
}

func GetVersion() string {
	out := MustExec(cmdGetVersion)

	// The first line is the zfs version, the second is the kernel module version
	version := strings.Split(out, "\n")[0]

	return version
}

func LoadKey(dataset string, passphrase string) (string, error) {
	cmd := append(cmdLoadKey, dataset)
	_, stderr, err := ExecWithInput(cmd, []byte(passphrase))

	return stderr, err
}

func UnloadKey(dataset string) (string, error) {
	cmd := append(cmdUnloadKey, dataset)
	_, stderr, err := Exec(cmd)

	return stderr, err
}

func Scrub(pool string) (string, error) {
	cmd := append(cmdScrub, pool)
	_, stderr, err := Exec(cmd)

	return stderr, err
}

func PauseScrub(pool string) (string, error) {
	cmd := append(cmdPauseScrub, pool)
	_, stderr, err := Exec(cmd)

	return stderr, err
}

func Mount(dataset string) (string, error) {
	cmd := append(cmdMount, dataset)
	_, stderr, err := Exec(cmd)

	return stderr, err
}

func Unmount(dataset string) (string, error) {
	cmd := append(cmdUnmount, dataset)
	_, stderr, err := Exec(cmd)

	return stderr, err
}

func Encode(raw interface{}) []byte {
	encoded, err := json.Marshal(raw)

	if err != nil {
		log.Fatalf("Unable to marshal %#v to JSON: %s", raw, err)
	}

	return encoded
}

func parseContainer(line string) structs.Container {
	info := strings.Fields(line)
	name := ""
	state := ""
	read := ""
	write := ""
	cksum := ""
	status := ""
	level := 0

	if len(info) == 0 {
		return structs.Container {
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

	return structs.Container {
		Name:   name,
		State:  state,
		Read:   read,
		Write:  write,
		Cksum:  cksum,
		Status: status,
		Level:  level,
	}
}
