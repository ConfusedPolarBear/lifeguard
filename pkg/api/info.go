// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"fmt"
	"net/http"

	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"
	"github.com/ConfusedPolarBear/lifeguard/pkg/config"

	"github.com/gorilla/mux"
)

func SetupInfo(r *mux.Router) {
	r.HandleFunc("/api/v0/info", infoHandler).Methods("GET")
	r.HandleFunc("/api/v0/support", supportHandler).Methods("GET")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	auth := checkSessionAuthQuiet(r)
	info := make(map[string]interface{})

	info["Product"] = "Lifeguard"
	info["Authenticated"] = auth
	info["Debug"] = config.DevMode

	if auth {
		info["ZFSVersion"] = zpool.GetVersion()

		info["Commit"] = config.Commit + config.Modified
		info["BuildTime"] = config.BuildTime
		info["GoVersion"] = config.GoVersion
	}

	EncodeAndSend(w, info)
}

func supportHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r, w)
	if username == "" {
		return
	}

	userAgent := r.UserAgent()
	zfsVersion := zpool.GetVersion()
	buildInfo := fmt.Sprintf("commit %s built at %s with %s", config.Commit + config.Modified, config.BuildTime, config.GoVersion)
	lsb := zpool.MustExec([]string { "/usr/bin/lsb_release", "-d"})
	kernel := zpool.MustExec([]string { "/bin/uname", "-r"})

	lsb = lsb[:len(lsb) - 1]

	response := fmt.Sprintf(`Lifeguard information:
	Username: %s
	Browser info: %s
	Build: %s
	
ZFS information:
	ZFS version: %s
	LSB %s
	Kernel: %s`,
	username, userAgent, buildInfo, zfsVersion, lsb, kernel)

	w.Write([]byte(response))
}