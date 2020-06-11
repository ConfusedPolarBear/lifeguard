// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"log"
	"net/http"
	"os"

	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const PORT string = ":5120"

var (
	key = []byte(os.Getenv("LIFEGUARD_KEY"))
	store = sessions.NewCookieStore(key)
	credentials = make(map[string]string)
)

func Setup() {
	// Validate session options
	if len(key) != 32 {
		log.Println("Set the LIFEGUARD_KEY environment variable to a securely generated random value and relaunch")
		log.Fatal("Invalid session key provided - it must be 32 bytes long")
	}

	// Validate temporary password
	temp := os.Getenv("LIFEGUARD_PASSWORD")
	if temp == "" {
		log.Println("Set the LIFEGUARD_PASSWORD environment variable to a bcrypt hash that will be used to login")
		log.Fatal("Invalid password hash provided")

	} else {
		log.Printf("Set password hash for user admin as '%s'", temp)
		credentials["admin"] = temp
	}

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,   // 30 days
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))

	// TODO: make these a tuple pair and set them up with a for-each loop
	http.HandleFunc("/api/v1/authenticate", loginHandler)
	http.HandleFunc("/api/v1/pools", getPoolsHandler)

	log.Printf("Listening on port %s, all interfaces", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func getPoolsHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSessionAuth(r, w) {
		return
	}

	w.Write(zpool.GetPools())
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	r.ParseForm()

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

	return session
}

func checkSessionAuth(r *http.Request, w http.ResponseWriter) (bool) {
	session := getSession(r)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Printf("%s cannot access %s: not authenticated", r.RemoteAddr, r.URL)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return false
	}

	return true
}
