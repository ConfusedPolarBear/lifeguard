// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"strings"
)

// These are set automatically by build.sh - DO NOT CHANGE THEM HERE
var Commit    string
var Modified  string
var GoVersion string

func Normalize() {
	GoVersion = strings.ReplaceAll(GoVersion, "go version ", "")
	if len(Commit) > 7 {
		Commit = Commit[:7]
	}
}