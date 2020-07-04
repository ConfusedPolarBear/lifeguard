// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"
	"strings"

	"github.com/ConfusedPolarBear/lifeguard/pkg/api"
	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Normalize generated config values
	config.GoVersion = strings.ReplaceAll(config.GoVersion, "go version ", "")
	if len(config.Commit) > 7 {
		config.Commit = config.Commit[:7]
	}

	log.Println("Starting Lifeguard")
	log.Printf("Git commit: %s%s", config.Commit, config.Modified)
	log.Printf("Build time: %s", config.BuildTime)
	log.Printf("Go version: %s", config.GoVersion)

	config.Load()

	api.Setup()
}
