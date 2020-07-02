// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package zpool

// These must be absolute paths to prevent executing arbitrary binaries
const cmdZpool = "/sbin/zpool"
const cmdZfs   = "/sbin/zfs"

// Basic information retrieval operations
var cmdGetVersion   = []string { cmdZpool, "version" }
var cmdPoolStatus   = []string { cmdZpool, "status" }
var cmdListDatasets = []string { cmdZfs, "list", "-p", "-H", "-o" }
var cmdListPools    = []string { cmdZpool, "list", "-p", "-H", "-o" }

// Cryptographic operations
var cmdLoadKey   = []string { cmdZfs, "load-key" }
var cmdUnloadKey = []string { cmdZfs, "unload-key" }

// Commands without zfs allow support
var cmdScrub      = []string { "sudo", cmdZpool, "scrub" }
var cmdPauseScrub = []string { "sudo", cmdZpool, "scrub", "-p" }
var cmdMount      = []string { "sudo", cmdZfs, "mount" }
var cmdUnmount    = []string { "sudo", cmdZfs, "unmount" }