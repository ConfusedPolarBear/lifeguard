// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"

	"github.com/ConfusedPolarBear/lifeguard/pkg/api"
	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.Normalize();

	log.Println("Starting Lifeguard")
	log.Printf("Git commit: %s%s", config.Commit, config.Modified)
	log.Printf("Build time: %s", config.BuildTime)
	log.Printf("Go version: %s", config.GoVersion)

	config.Load()

	api.Setup()
}
