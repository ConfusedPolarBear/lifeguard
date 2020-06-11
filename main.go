// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"

	"github.com/ConfusedPolarBear/lifeguard/pkg/api"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	api.Setup()
}
