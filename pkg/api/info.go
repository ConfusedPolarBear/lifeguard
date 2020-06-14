// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"net/http"

	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"
)

func SetupInfo() {
	http.HandleFunc("/api/v0/info", versionHandler)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	auth := checkSessionAuthQuiet(r, w)
	info := make(map[string]interface{})

	info["Product"] = "Lifeguard"
	info["Authenticated"] = auth

	if auth {
		info["ZFSVersion"] = zpool.GetVersion()
	}

	w.Write(zpool.Encode(info))
}
