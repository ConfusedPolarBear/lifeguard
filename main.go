// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"
	"fmt"

	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	pools := string(zpool.GetPools())

	fmt.Println(pools)
}
