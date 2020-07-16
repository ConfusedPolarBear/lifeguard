// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"net/http"
	"log"

	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"
	
	"github.com/gorilla/mux"
)

func SetupTOTP(r *mux.Router) {
	r.HandleFunc("/api/v0/tfa/totp/initialize", initializeHandler).Methods("GET")
	r.HandleFunc("/api/v0/tfa/totp/save", saveHandler).Methods("POST")
	r.HandleFunc("/api/v0/tfa/totp/authenticate", challengeHandler).Methods("POST")
}

func initializeHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsernameQuiet(r)
	if username == "" {
		return
	}

	secret, image := config.InitializeTOTP(username)
	image = "data:image/png;base64," + image

	ret := struct {
		Secret string
		Image string
	} {
		secret,
		image,
	}

	w.Write(zpool.Encode(ret))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsernameQuiet(r)
	if username == "" {
		return
	}

	secret, secretOk := GetParameter(r, "secret")
	if !secretOk {
		ReportMissing(w)
		return
	}

	code, codeOk := GetParameter(r, "code")
	if !codeOk {
		ReportMissing(w)
		return
	}

	ok := config.SaveTOTP(username, secret, code)
	if !ok {
		http.Error(w, "Invalid code - check that the time is in sync on the server and your phone", http.StatusBadRequest)
		return
	}

	http.Error(w, "OK", http.StatusOK)
}

func challengeHandler(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	username := getPartialAuth(r)
	if username == "" {
		return
	}

	code, codeOk := GetParameter(r, "code")
	if !codeOk {
		ReportMissing(w)
		return
	}

	ok := config.VerifyTOTP(username, code)
	if !ok {
		http.Error(w, "Invalid code", http.StatusForbidden)
		return
	}

	session.Values["authenticated"] = true
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Unable to save session: %s", err)
		return
	}

	http.Error(w, "OK", http.StatusOK)
}