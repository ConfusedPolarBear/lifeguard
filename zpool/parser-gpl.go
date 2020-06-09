// Copyright 2020 Nathan Skelton
// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

// Structs originally from https://github.com/skeltonn/zfs-manager

package zpool

// Struct for "zpool status"
type PoolStatus struct {
	Name       string
	State      string
	Status     string
	Scan       string
	Action     string
	See        string
	Containers []*ContainerStatus
	Errors     string
	Raw        string
}

// Struct for within "zpool status", the config section
type ContainerStatus struct {
	Name   string
	State  string
	Read   string
	Write  string
	Cksum  string
	Status string
}
