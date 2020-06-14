// Copyright 2020 Nathan Skelton
// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

// Structs originally from https://github.com/skeltonn/zfs-manager

package zpool

// Pool wide status
type Pool struct {
	Name       string
	State      string
	Status     string
	Scan       string
	Action     string
	See        string
	Containers []*Container
	Errors     string
	Raw        string
	Properties map[string]*Property
}

// VDEV or VDEV member status
type Container struct {
	Name   string
	State  string
	Read   string
	Write  string
	Cksum  string
	Status string
	Level  int
}
