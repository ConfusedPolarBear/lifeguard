// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

// TODO: use something like cron to parse all pools every 15 seconds

package notifications

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ConfusedPolarBear/lifeguard/pkg/structs"
)

var Notifications []structs.Notification

func UpdatePoolState(pool string, current *structs.Pool, previous *structs.Pool) {
	// If this is the first update, there won't be a previous pool state to compare against
	if previous == nil {
		return
	}

	/* Notifications to implement:
	 *     Pool container read/write/checksum count change
	 *     Pool container state change
	*/

	// TODO: telegram library: https://github.com/tucnak/telebot

	// Notification 1: Pool wide state change
	if previous.State != current.State {
		SendNotification(1, "critical", fmt.Sprintf("Pool \"%s\" state changed: %s -> %s", pool, previous.State, current.State))
	}

	// Notification 2: Pool status change
	if previous.Status != current.Status {
		SendNotification(2, "warning", fmt.Sprintf("Pool \"%s\" new status: %s", pool, CleanupString(current.Status)))
	}

	// Notification 3: Errors change
	if previous.Errors != current.Errors {
		SendNotification(3, "critical", fmt.Sprintf("Pool \"%s\" new errors: %s", pool, current.Errors))
	}

	// Notification 4: Scrub start
	// When a scrub starts, the Scanned property will change from 0 (no scrub currently active) to a number that is not 0
	if previous.Scanned == 0 && current.Scanned != 0 {
		SendNotification(4, "info", fmt.Sprintf("Pool \"%s\" scrub: started", pool))
	}

	// Notification 5: Scrub finish
	// When a scrub finishes, the Scanned property will change from not 0 (currently scrubbing) to 0
	if previous.Scanned != 0 && current.Scanned == 0 {
		// TODO: parse and include how much was resilvered and any errors
		SendNotification(5, "info", fmt.Sprintf("Pool \"%s\" scrub: completed", pool))
	}
}

func SendNotification(id int, severity string, message string) {
	n := structs.Notification {
		ID: id,
		Timestamp: time.Now(),
		Severity: severity,
		Message: message,
	}

	log.Printf("Got notification %s", n.String())
}

func CleanupString(raw string) string {
	// TODO: replace with regex
	raw = strings.ReplaceAll(raw, "\r", " ")
	raw = strings.ReplaceAll(raw, "\n", " ")
	raw = strings.ReplaceAll(raw, "\t", " ")

	return raw
}