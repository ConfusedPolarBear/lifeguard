// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package zpool

const cmdZpool = "/sbin/zpool"
const cmdZfs   = "/sbin/zfs"

// pool.go
// TODO: use the -p flag to enforce exact (parsable) numbers in output for sizes?
var cmdPoolStatus   = []string { cmdZpool, "status" }
var cmdListPools    = []string { cmdZpool, "list", "-p", "-H", "-o" }
var cmdListDatasets = []string { cmdZfs, "list", "-p", "-H", "-o" }
var cmdGetVersion   = []string { cmdZpool, "version" }
