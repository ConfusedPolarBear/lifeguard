// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package zpool

// These must be absolute paths to prevent executing arbitrary binaries
const cmdZpool = "/sbin/zpool"
const cmdZfs   = "/sbin/zfs"

// pool.go
var cmdGetVersion   = []string { cmdZpool, "version" }
var cmdListDatasets = []string { cmdZfs, "list", "-p", "-H", "-o" }
var cmdListPools    = []string { cmdZpool, "list", "-p", "-H", "-o" }
var cmdLoadKey      = []string { cmdZfs, "load-key" }
var cmdPoolStatus   = []string { cmdZpool, "status" }

// pool.go commands without zfs allow support
var cmdScrub        = []string { "sudo", cmdZpool, "scrub" }
