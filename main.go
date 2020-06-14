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

	config.Load()

	api.Setup()
}
