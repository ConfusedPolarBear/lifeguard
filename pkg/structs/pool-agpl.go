// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package structs

type Property struct {
	Name  string
	Value string
	HMAC  string
}

// Reused to represent datasets and snapshots
type Data struct {
	Name       string
	Type       string
	Properties map[string]*Property
	Internal   map[string]*Property
}
