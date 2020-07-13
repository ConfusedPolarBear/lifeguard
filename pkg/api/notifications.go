// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"net/http"

	"github.com/ConfusedPolarBear/lifeguard/pkg/notifications"
	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"

	"github.com/gorilla/mux"
)

func SetupNotifications(r *mux.Router) {
	notifications.Initialize()

	// List all notifications
	r.HandleFunc("/api/v0/notifications/list", getNotifications).Methods("GET")
}

func getNotifications(w http.ResponseWriter, r *http.Request) {
	if !checkSessionAuth(r, w) {
		return
	}

	w.Write(zpool.Encode(notifications.Notifications))
}