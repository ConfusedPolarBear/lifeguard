// Copyright 2020 Nathan Skelton
// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

// Structs originally from https://github.com/skeltonn/zfs-manager

package structs

// Pool wide status
type Pool struct {
	Name       string
	State      string
	Status     string
	Scan       string
	Scanned    float64
	ScanPaused bool
	Action     string
	See        string
	Containers []*Container
	Errors     string
	Raw        string
	Datasets   []map[string]*Property
	Snapshots  []map[string]*Property
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

var DefaultProperties = map[string]string {
	"pool":     "name,health,capacity,free,size,fragmentation,ashift",
	"dataset":  "name,used,avail,encryption,keystatus,mounted,usedsnap,usedds",
	"snapshot": "name,used,avail,refer",
}