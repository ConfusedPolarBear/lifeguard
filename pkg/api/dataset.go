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
)

func SetupDataset() {
	// Retrieves information for a dataset or snapshot
	http.HandleFunc("/api/v0/data", getDataInfoHandler)
	http.HandleFunc("/api/v0/data/mount", mountHandler)
	http.HandleFunc("/api/v0/data/unmount", unmountHandler)

	// Load and unload encryption keys
	http.HandleFunc("/api/v0/key/load", loadKeyHandler)
	http.HandleFunc("/api/v0/key/unload", unloadKeyHandler)

	// Start or pause a pool scrub
	http.HandleFunc("/api/v0/pool/scrub", scrubHandler)
	http.HandleFunc("/api/v0/pool/scrub/pause", scrubPauseHandler)
}

func getDataInfoHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSessionAuth(r, w) {
		return
	}

	var saved *zpool.Data

	name, ok := GetHMAC(w, r)
	if !ok {
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

	passphrase, okPassphrase := GetParameter(w, r, "passphrase")
	if !okPassphrase {
		return
	}

	name, okName := GetHMAC(w, r)
	if !okName {
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

	name, okName := GetHMAC(w, r)
	if !okName {
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
	
	name, ok := GetHMAC(w, r)
	if !ok {
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
	
	name, ok := GetHMAC(w, r)
	if !ok {
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
	
	name, ok := GetHMAC(w, r)
	if !ok {
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
	
	name, ok := GetHMAC(w, r)
	if !ok {
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