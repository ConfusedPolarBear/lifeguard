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
	"github.com/ConfusedPolarBear/lifeguard/pkg/structs"
	"github.com/ConfusedPolarBear/lifeguard/pkg/zpool"

	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	key = []byte("")			// use a temporary key so key and store are accessible throughout the api package
	store = sessions.NewCookieStore(key)
)

// This is used by getPropertiesHandler to construct the fields object. The custom JSON fields are needed because go won't export struct members with a lowercase name.
// However, bootstrap vue requires the fields to be in dromedary case (first letter lowercase)
type Column struct {
	Key string		`json:"key"`
	Sortable bool	`json:"sortable"`
}

func Setup() {
	port := config.GetString("bind", "")
	if port == "" {
		log.Printf("Warning: No option was specified for server bind address, listening on port 5120 (all interfaces)")
		port = ":5120"
	}

	key = []byte(config.GetString("keys.session", ""))

	// Validate session options
	if len(key) != 32 {
		log.Println("Regenerating session key")

		temp := strings.ToUpper(crypto.GetRandom(16))
		key = []byte(temp)
		config.Set("keys.session", temp)
	}

	store = sessions.NewCookieStore(key)

	if !config.IsUser("admin") {
		fmt.Print("Enter new password for user admin: ")

		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Unable to get password: %s", err)
		}

		config.CreateUser("admin", crypto.HashPassword(string(bytePassword)), nil)

		fmt.Println()
		log.Printf("Password successfully hashed and saved")
	}

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
	r.HandleFunc("/api/v0/tfa/enabled", checkTFAEnabledHandler).Methods("GET")
	r.HandleFunc("/api/v0/tfa/challenge", tfaChallengeHandler).Methods("GET")

	// Pool
	r.HandleFunc("/api/v0/pool/{pool}", getPoolHandler).Methods("GET")
	r.HandleFunc("/api/v0/pools", getAllPoolsHandler).Methods("GET")
	r.HandleFunc("/api/v0/properties/{type}", getPropertyListHandler).Methods("GET")

	SetupInfo(r)
	SetupDataset(r)
	SetupNotifications(r)
	SetupTOTP(r)

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
		// TODO: allow unsafe-eval in dev? disabling eval breaks vue dev tools in firefox 78.0.1
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
		props = config.GetString("properties.dataset", structs.DefaultProperties["dataset"])
	} else if name == "Snapshots" {
		props = config.GetString("properties.snapshot", structs.DefaultProperties["snapshot"])
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
	partialAuth := "full"
	
	session.Values["authenticated"] = auth
	session.Values["username"] = username

	if auth && config.IsTwoFactorEnabled(username) {
		partialAuth = "partial"
		session.Values["partialAuth"] = username	
		session.Values["authenticated"] = false
	}

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Unable to save session: %s", err)
		return
	}

	if auth {
		log.Printf("%s authenticated (%s) from %s", username, partialAuth, r.RemoteAddr)
		http.Error(w, partialAuth, http.StatusOK)

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

func tfaChallengeHandler(w http.ResponseWriter, r *http.Request) {
	username := getPartialAuth(r)
	if username == "" {
		http.Error(w, "Invalid state", http.StatusForbidden)
		return
	}

	provider := config.GetTwoFactorProvider(username)
	challenge := ""

	// challenge will be needed for fido2
	if provider == "totp" {
		challenge = ""
	}

	ret := struct {
		Provider string
		Challenge string
	} {
		provider,
		challenge,
	}

	w.Write(zpool.Encode(ret))
}

func getAuth(r *http.Request) (string, string) {
	username, _ := GetParameter(r, "Username")
	password, _ := GetParameter(r, "Password")

	return username, password
}

func checkAuth(username string, password string) (bool, string) {
	goodUsername := true
	user, ok := config.GetUsers()[username]

	if !ok {
		log.Printf("Unknown username %s", username)

		// Prevent user enumeration attacks
		user.Password = "$2a$12$000000000000.0000000000000000000000000000000000000000"
		goodUsername = false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

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

func getPartialAuth(r *http.Request) string {
	session := getSession(r)
	partial, ok := "", false

	if partial, ok = session.Values["partialAuth"].(string); !ok {
		return ""
	}

	return partial
}

func checkTFAEnabledHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsernameQuiet(r)
	if username == "" {
		return
	}

	enabled := config.IsTwoFactorEnabled(username)

	ret := struct {
		Enabled bool
	} {
		enabled,
	}

	w.Write(zpool.Encode(ret))
}