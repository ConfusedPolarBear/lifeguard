// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

// TODO: add an installation note about permissions that need to be given to the user the program runs as

package api

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"path/filepath"

	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
	"github.com/ConfusedPolarBear/lifeguard/pkg/crypto"
	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"

	"github.com/gorilla/mux"
)

type File struct {
	Type string
	Name string
	HMAC string
	Size string
}

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

	// File browsing
	r.HandleFunc("/api/v0/files/browse/{id}", browseFilesHandler).Methods("GET")		// list directory
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

func browseFilesHandler(w http.ResponseWriter, r *http.Request) {
	var files []File
	username := getUsername(r, w)
	if username == "" {
		log.Printf("not authed")
		return
	}

	path, ok := GetHMAC(r)
	if !ok {
		log.Printf("bad hmac")
		return
	}

	// Dataset names aren't prefixed with a slash
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	log.Printf("%s browsed to %s", username, path)

	// TODO: Lifeguard could verify that the browser binary has a signature on it
	// The signature can be from any key but the public key must be printed at startup
	//    and must be the same between all binaries
	contents, stderr, err := zpool.Exec([]string { "./browser", "-f", path })

	if len(contents) <= 3 || err != nil {
		http.Error(w, "Unable to list contents", http.StatusInternalServerError)
		log.Printf("Unable to list contents of %s: %s. Error: %s", path, err, stderr)
		return
	}

	// test if this is a folder or file
	if contents[:4] == "fold" {
		contents = contents[4:]

		// The first item with type "@" is the current path that we are at
	files = append(files, File {
		Type: "@",
		Name: path,
		HMAC: "",
		Size: "0",
	})

	for _, raw := range strings.Split(contents, "\n") {
		parts := strings.Split(raw, " ")
		if len(parts) != 3 {
			break
		}

		decoded, _ := base64.StdEncoding.DecodeString(parts[1])
		
		current := path + "/" + string(decoded)
		hmac := crypto.GenerateHMAC(current)

		files = append(files, File {
			Type: parts[0],
			Name: string(decoded),
			HMAC: hmac,
			Size: parts[2],
		})
	}

	w.Write(zpool.Encode(files))

	} else if contents[:4] == "file" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(path)))
		w.Write([]byte(contents[4:]))

	} else {
		http.Error(w, "Unknown type", http.StatusInternalServerError)
	}
}