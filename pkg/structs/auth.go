// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package structs

type User struct {
	Username          string
	Password          string
	TwoFactorProvider string
	TwoFactorData     string
}