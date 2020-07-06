// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

// TODO: add an installation note about permissions that need to be given to the user the program runs as

package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"

	"github.com/gorilla/mux"
)

func SetupDataset(r *mux.Router) {
	// Retrieves information for a dataset or snapshot
	r.HandleFunc("/api/v0/data/{id}/info", getDataInfoHandler).Methods("GET")
	r.HandleFunc("/api/v0/data/{id}/mount", mountHandler).Methods("POST")
	r.HandleFunc("/api/v0/data/{id}/unmount", unmountHandler).Methods("POST")

	// Load and unload encryption keys
	r.HandleFunc("/api/v0/key/{id}/load", loadKeyHandler).Methods("POST")
	r.HandleFunc("/api/v0/key/{id}/unload", unloadKeyHandler).Methods("POST")

	// Start or pause a pool scrub
	r.HandleFunc("/api/v0/pool/{id}/scrub/start", scrubHandler).Methods("POST")
	r.HandleFunc("/api/v0/pool/{id}/scrub/pause", scrubPauseHandler).Methods("POST")
}

func getDataInfoHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSessionAuth(r, w) {
		return
	}

	var saved *zpool.Data

	name, ok := GetHMAC(r)
	if !ok {
		ReportInvalid(w)
		return
	}

	// Grab the keylocation so the load key button can be conditionally enabled
	saved = &zpool.Data {
		Name:       name,
		Type:       zpool.GetProperties(name, "zfs", "", "type")[0]["type"].Value,
		Properties: zpool.GetProperties(name, "zfs", "", config.GetString("properties.dataset"))[0],
		Internal:   zpool.GetProperties(name, "zfs", "", "keylocation")[0],
	}

	w.Write(zpool.Encode(saved))
}

func loadKeyHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r, w)
	if username == "" {
		return
	}

	passphrase, okPassphrase := GetParameter(r, "Passphrase")
	name, okName := GetHMAC(r)

	if !okPassphrase || !okName {
		ReportMissing(w)
		return
	}

	stderr, err := zpool.LoadKey(name, passphrase)
	if err != nil {
		// TODO: unit test the first two conditions
		if strings.Index(stderr, "Incorrect key provided") != -1 {
			http.Error(w, "Incorrect passphrase", http.StatusUnauthorized)

		} else if strings.Index(stderr, "Key already loaded") != -1 {
			http.Error(w, "", http.StatusOK)

		} else {
			log.Printf("Unable to load key for %s: %s. %s", name, err, stderr)
			http.Error(w, msgErrorOccurred, http.StatusBadRequest)
		}

		return
	}

	log.Println(fmt.Sprintf("%s loaded key for %s", username, name))
	http.Error(w, "", http.StatusOK)
}

func unloadKeyHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r, w)
	if username == "" {
		return
	}

	name, okName := GetHMAC(r)
	if !okName {
		ReportInvalid(w)
		return
	}

	stderr, err := zpool.UnloadKey(name)
	if err != nil {
		if strings.Index(stderr, "is busy") != -1 {
			http.Error(w, "Dataset is mounted", http.StatusBadRequest)

		} else {
			log.Printf("Unable to unload key for %s: %s. %s", name, err, stderr)
			http.Error(w, msgErrorOccurred, http.StatusBadRequest)
		}

		return
	}

	log.Println(fmt.Sprintf("%s unloaded key for %s", username, name))
	http.Error(w, "", http.StatusOK)
}

func scrubHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r, w)
	if username == "" {
		return
	}
	
	name, ok := GetHMAC(r)
	if !ok {
		ReportInvalid(w)
		return
	}

	stderr, err := zpool.Scrub(name)

	if err != nil {
		log.Printf("Unable to scrub pool %s: %s. %s", name, err, stderr)
		http.Error(w, stderr, http.StatusBadRequest)

		return
	}

	log.Println(fmt.Sprintf("%s started scrub for pool %s", username, name))
	http.Error(w, "", http.StatusOK)
}

// TODO: dedup this code with scrubHandler?
func scrubPauseHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r, w)
	if username == "" {
		return
	}
	
	name, ok := GetHMAC(r)
	if !ok {
		ReportInvalid(w)
		return
	}

	stderr, err := zpool.PauseScrub(name)

	if err != nil {
		log.Printf("Unable to pause scrubbing pool %s: %s. %s", name, err, stderr)
		http.Error(w, stderr, http.StatusBadRequest)

		return
	}

	log.Println(fmt.Sprintf("%s paused scrub for pool %s", username, name))
	http.Error(w, "", http.StatusOK)
}

func mountHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r, w)
	if username == "" {
		return
	}
	
	name, ok := GetHMAC(r)
	if !ok {
		ReportInvalid(w)
		return
	}

	stderr, err := zpool.Mount(name)

	if err != nil {
		if strings.Index(stderr, "encryption key not loaded") != -1 {
			http.Error(w, "Encryption key is not loaded", http.StatusBadRequest)
		} else {
			log.Printf("Unable to mount dataset %s: %s. %s", name, err, stderr)
			http.Error(w, msgErrorOccurred, http.StatusBadRequest)
		}

		return
	}

	log.Println(fmt.Sprintf("%s mounted dataset %s", username, name))
	http.Error(w, "", http.StatusOK)
}

func unmountHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r, w)
	if username == "" {
		return
	}
	
	name, ok := GetHMAC(r)
	if !ok {
		ReportInvalid(w)
		return
	}

	stderr, err := zpool.Unmount(name)

	if err != nil {
		log.Printf("Unable to unmount dataset %s: %s. %s", name, err, stderr)
		http.Error(w, msgErrorOccurred, http.StatusBadRequest)

		return
	}

	log.Println(fmt.Sprintf("%s unmounted dataset %s", username, name))
	http.Error(w, "", http.StatusOK)
}