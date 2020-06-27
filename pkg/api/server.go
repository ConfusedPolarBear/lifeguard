// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

// TODO: save auth info in the session (username, 2fa, etc)

package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"syscall"

	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
	"github.com/ConfusedPolarBear/lifeguard/pkg/crypto"
	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

const PORT string = ":5120"
const FORM_SIZE int64 = 2048

var (
	key = []byte("")			// use a temporary key so key and store are accessible throughout the api package
	store = sessions.NewCookieStore(key)
	credentials = make(map[string]string)
)

func Setup() {
	key = []byte(config.GetString("security.session_key"))

	// Validate session options
	if len(key) != 32 {
		log.Println("Regenerating session key")

		temp := strings.ToUpper(crypto.GetRandom(16))
		key = []byte(temp)
		config.Set("security.session_key", temp)
	}

	store = sessions.NewCookieStore(key)

	adminHash := config.GetString("security.admin")
	if adminHash == "" {
		fmt.Print("Enter new password for user admin: ")

		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Unable to get password: %s", err)
		}

		adminHash = crypto.HashPassword(string(bytePassword))
		config.Set("security.admin", adminHash)

		fmt.Println()
		log.Printf("Password successfully hashed and saved")
	}

	credentials["admin"] = adminHash

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,   // 30 days
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}

	// Static web UI
	http.Handle("/", http.FileServer(http.Dir("./web/dist")))

	// Security
	http.HandleFunc("/api/v0/authenticate", loginHandler)
	http.HandleFunc("/api/v0/logout", logoutHandler)

	// Pool
	http.HandleFunc("/api/v0/pool", getPoolHandler)
	http.HandleFunc("/api/v0/pools", getAllPoolsHandler)

	SetupInfo()
	SetupDataset()

	log.Printf("Listening on port %s, all interfaces", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func GetParameter(w http.ResponseWriter, r *http.Request, name string) (string, bool) {
	data := ""
	if rawData, ok := r.Form[name]; ok {
		data = rawData[0]
	} else {
		msg := fmt.Sprintf("Missing %s parameter", name)
		http.Error(w, msg, http.StatusBadRequest)
		return "", false
	}

	return data, true
}

func GetHMAC(w http.ResponseWriter, r *http.Request) (string, bool) {
	data := ""
	hmac, ok := GetParameter(w, r, "id")
	if ok {
		data = crypto.LookupHMAC(hmac)
		if data == "" {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return "", false
		}
	} else {
		return "", false
	}

	return data, true
}

func getAllPoolsHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSessionAuth(r, w) {
		return
	}

	pools := zpool.ParseAllPools()
	w.Write(zpool.Encode(pools))
}

func getPoolHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSessionAuth(r, w) {
		return
	}

	pool, ok := GetParameter(w, r, "pool")
	if !ok {
		return
	}

	parsed := zpool.ParsePool(pool, true)
	w.Write(zpool.Encode(parsed))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	username, password := getAuth(r)
	auth := checkAuth(username, password)
	session.Values["authenticated"] = auth

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Unable to save session: %s", err)
		return
	}

	if auth {
		log.Printf("%s authenticated from %s", username, r.RemoteAddr)
		http.Error(w, "OK", http.StatusOK)

	} else {
		log.Printf("WARNING %s failed to authenticate as %s", r.RemoteAddr, username)
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	session.Values = nil

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Unable to save session: %s", err)
		return
	}

	http.Error(w, "", http.StatusOK)
}

func getAuth(r *http.Request) (string, string) {
	username := ""
	password := ""

	if rawUsername, ok := r.Form["Username"]; ok {
		username = rawUsername[0]
	}

	if rawPassword, ok := r.Form["Password"]; ok {
		password = rawPassword[0]
	}

	return username, password
}

func checkAuth(username string, password string) (bool) {
	goodUsername := true
	knownGood, ok := credentials[username]

	if !ok {
		log.Printf("Unknown username %s", username)

		// Prevent user enumeration attacks
		knownGood = "$2a$12$000000000000.0000000000000000000000000000000000000000"
		goodUsername = false
	}

	err := bcrypt.CompareHashAndPassword([]byte(knownGood), []byte(password))

	if err != nil {
		if goodUsername {
			log.Printf("Incorect password provided for %s", username)
		}

		return false
	}

	return true
}

func getSession(r *http.Request) (*sessions.Session) {
	// store.Get returns a blank session if there is an error, so the error is safe to ignore
	session, _ := store.Get(r, "SESSION")

	r.ParseMultipartForm(FORM_SIZE)

	return session
}

func checkSessionAuthInternal(r *http.Request, w http.ResponseWriter) (bool) {
	session := getSession(r)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Printf("%s cannot access %s: not authenticated", r.RemoteAddr, r.URL)
		return false
	}

	return true
}

func checkSessionAuthQuiet(r *http.Request, w http.ResponseWriter) (bool) {
	return checkSessionAuthInternal(r, w)
}

func checkSessionAuth(r *http.Request, w http.ResponseWriter) (bool) {
	auth := checkSessionAuthInternal(r, w)

	if !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	return auth
}
