// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package zpool

const cmdZpool = "/sbin/zpool"
const cmdZfs = "/sbin/zfs"

// pool.go
var cmdPoolStatus   = []string { cmdZpool, "status" }
var cmdListPools    = []string { cmdZpool, "list", "-H", "-o" }
var cmdListDatasets = []string { cmdZfs, "list", "-H", "-o" }
