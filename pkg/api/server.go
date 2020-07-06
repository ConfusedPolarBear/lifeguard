// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
	"github.com/ConfusedPolarBear/lifeguard/pkg/crypto"
	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"

	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	key = []byte("")			// use a temporary key so key and store are accessible throughout the api package
	store = sessions.NewCookieStore(key)
	credentials = make(map[string]string)
)

// This is used by getPropertiesHandler to construct the fields object. The custom JSON fields are needed because go won't export struct members with a lowercase name.
// However, bootstrap vue requires the fields to be in dromedary case (first letter lowercase)
type Column struct {
	Key string		`json:"key"`
	Sortable bool	`json:"sortable"`
}

func Setup() {
	port := config.GetString("server.bind")
	if port == "" {
		log.Printf("Warning: No option was specified for server.bind, listening on port 5120 (all interfaces)")
		port = ":5120"
	}

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
		MaxAge:   30 * 86400,   // 30 days
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}

	r := mux.NewRouter()

	// Security
	r.HandleFunc("/api/v0/authenticate", loginHandler).Methods("POST")
	r.HandleFunc("/api/v0/logout", logoutHandler).Methods("POST")

	// Pool
	r.HandleFunc("/api/v0/pool/{pool}", getPoolHandler).Methods("GET")
	r.HandleFunc("/api/v0/pools", getAllPoolsHandler).Methods("GET")
	r.HandleFunc("/api/v0/properties/{type}", getPropertyListHandler).Methods("GET")

	SetupInfo(r)
	SetupDataset(r)

	// Static web UI
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/dist"))).Methods("GET")

	// Middleware
	r.Use(securityHeadersMw)

	srv := &http.Server{
		Handler:      r,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on %s", port)
	log.Fatal(srv.ListenAndServe())
}

func securityHeadersMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: eliminate data URLs for images
		w.Header().Set("Content-Security-Policy", `default-src 'self';
			base-uri 'none';
			block-all-mixed-content;
			form-action 'self';
			frame-ancestors 'none';
			img-src 'self' data:;
			object-src 'none';`)

		w.Header().Set("Referrer-Policy", "no-referrer")		// never send referrer
		w.Header().Set("X-Frame-Options", "deny")				// forbid framing
		w.Header().Set("X-Content-Type-Options", "nosniff")		// forbid content type sniffing

		next.ServeHTTP(w, r)
	})
}

func GetParameter(r *http.Request, name string) (string, bool) {
	// First search for the variable in the URL
	vars := mux.Vars(r)
	if rawVar, ok := vars[name]; ok {
		return string(rawVar), true
	}

	// If the variable wasn't in the URL, it is encoded as a URL encoded form
	r.ParseForm()
	data := r.FormValue(name)

	return data, (data != "")
}

func ReportMissing(w http.ResponseWriter) {
	msg := fmt.Sprintf("Missing parameter")
	http.Error(w, msg, http.StatusBadRequest)
}

func ReportInvalid(w http.ResponseWriter) {
	msg := fmt.Sprintf("Invalid id")
	http.Error(w, msg, http.StatusBadRequest)
}

func GetHMAC(r *http.Request) (string, bool) {
	data := ""
	hmac, ok := GetParameter(r, "id")
	if ok {
		data = crypto.LookupHMAC(hmac)
		if data == "" {
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

	pool, ok := GetParameter(r, "pool")
	if !ok {
		ReportMissing(w)
		return
	}

	parsed := zpool.ParsePool(pool, true)
	w.Write(zpool.Encode(parsed))
}

// This handler returns the properties the user has specified in the config.ini file to sort the displayed columns correctly
func getPropertyListHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSessionAuth(r, w) {
		return
	}

	name, ok := GetParameter(r, "type")
	if !ok {
		ReportMissing(w)
		return
	}

	props := ""
	if name == "Datasets" {
		props = config.GetString("properties.dataset")
	} else if name == "Snapshots" {
		props = config.GetString("properties.snapshot")
	} else {
		http.Error(w, "Unknown value for type parameter", http.StatusBadRequest)
		return
	}

	var columns []*Column

	// TODO: why is sorting by name broken?
	for _, col := range strings.Split(props, ",") {
		columns = append(columns, &Column {
			Key: col,
			Sortable: (col != "name"),
		})
	}

	w.Write(zpool.Encode(columns))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	sentUsername, password := getAuth(r)
	auth, username := checkAuth(sentUsername, password)
	
	session.Values["authenticated"] = auth
	session.Values["username"] = username

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
		log.Printf("WARNING %s failed to authenticate as %s", r.RemoteAddr, sentUsername)
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
	username, _ := GetParameter(r, "Username")
	password, _ := GetParameter(r, "Password")

	return username, password
}

func checkAuth(username string, password string) (bool, string) {
	goodUsername := true
	loadedHash, ok := credentials[username]

	if !ok {
		log.Printf("Unknown username %s", username)

		// Prevent user enumeration attacks
		loadedHash = "$2a$12$000000000000.0000000000000000000000000000000000000000"
		goodUsername = false
	}

	err := bcrypt.CompareHashAndPassword([]byte(loadedHash), []byte(password))

	if err != nil {
		if goodUsername {
			log.Printf("Incorect password provided for %s", username)
		}

		return false, ""
	}

	return true, username
}

func getSession(r *http.Request) (*sessions.Session) {
	// store.Get returns a blank session if there is an error, so the error is safe to ignore
	session, _ := store.Get(r, "SESSION")

	return session
}

func checkSessionAuthInternal(r *http.Request) (bool) {
	session := getSession(r)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Printf("%s cannot access %s: not authenticated", r.RemoteAddr, r.URL)
		return false
	}

	return true
}

func checkSessionAuthQuiet(r *http.Request) (bool) {
	return checkSessionAuthInternal(r)
}

func checkSessionAuth(r *http.Request, w http.ResponseWriter) (bool) {
	auth := checkSessionAuthInternal(r)

	if !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	return auth
}

func getUsernameInternal(r *http.Request) (string) {
	if !checkSessionAuthQuiet(r) {
		return ""
	}

	return getSession(r).Values["username"].(string)
}

func getUsernameQuiet(r *http.Request) string {
	return getUsernameInternal(r)
}

func getUsername(r *http.Request, w http.ResponseWriter) string {
	username := getUsernameInternal(r)

	if username == "" {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	return username
}