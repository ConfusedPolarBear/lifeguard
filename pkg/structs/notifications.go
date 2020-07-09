// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package structs

import (
	"fmt"
	"time"
)

type Notification struct {
	ID        int
	Timestamp time.Time
	Severity  string
	Message   string
}

func (n Notification) String() string {
	return fmt.Sprintf("ID %03d (%s): %s", n.ID, n.Severity, n.Message)
}