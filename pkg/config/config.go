// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"log"

	"github.com/ConfusedPolarBear/lifeguard/pkg/structs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var db *sql.DB
var driver = "sqlite3"
var connString = "./config/config.db"

func Load() {
	err := errors.New("OK")

	db, err = sql.Open(driver, connString)
	if errors.Is(err, errors.New("OK")) {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	prepare("create table if not exists config (Key string primary key unique, Value string not null)").Exec()
	prepare("create table if not exists auth (Username string primary key unique, Password string not null, TwoFactor string)").Exec()

	if loadLegacy() {
		path := viper.ConfigFileUsed()

		log.Printf("Migrating legacy configuration %s to database", path)

		loadLegacy()
		migrateFromLegacy()

		dst := path + ".bak"
		os.Rename(path, dst)
		log.Printf("Legacy configuration renamed: %s -> %s", path, dst)
	}
}

func prepare(raw string) *sql.Stmt {
	stmt, err := db.Prepare(raw)
	if err != nil {
		log.Fatalf("Unable to prepare statement %s: %s", raw, err)
	}

	return stmt
}

func loadLegacy() bool {
	viper.SetConfigName("config")

	viper.AddConfigPath("./config/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return false
	}

	return true
}

// TODO: delete this function after all users have moved over
func migrateFromLegacy() {
	// Get the number of row in the config database - if it's zero, all values are still in the legacy config.ini file
	var count int
	row, _ := db.Query("SELECT count(*) from config")
	row.Next()
	row.Scan(&count)
	row.Close()
	if count != 0 {
		return
	}
	
	// Use a transaction to abort the migration if something fails
	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Migration failed: unable to create migration transaction: %s", err)
	}

	// Stage 1: Migrate settings
	migrateConfig("exec.timeout", tx)
	migrateConfig("exec.timeout_path", tx)
	migrateConfig("properties.dataset", tx)
	migrateConfig("properties.pool", tx)
	migrateConfig("properties.snapshot", tx)
	migrateConfig("security.session_key", tx)
	migrateConfig("server.bind", tx)
	migrateConfig("debug.exec", tx)
	migrateConfig("debug.parse", tx)
	
	// Stage 2: Migrate the admin account
	CreateUser("admin", viper.GetString("security.admin"), tx)

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Migration failed: unable to commit migration transaction: %s", err)
	}

	log.Printf("Migration completed successfully")
}

func migrateConfig(key string, tx *sql.Tx) {
	var mappings = make(map[string]string)
	mappings["exec.timeout"] = "timeout.value"
	mappings["exec.timeout_path"] = "timeout.path"
	mappings["security.session_key"] = "keys.session"
	mappings["server.bind"] = "bind"

	value := viper.GetString(key)
	if value == "" {
		log.Fatalf("Migration failed for key %s: unable to find key", key)
	}

	if newKey, ok := mappings[key]; ok {
		key = newKey
	}
	
	stmt := tx.Stmt(prepare("insert into config values (?, ?)"))
	_, err := stmt.Exec(key, value)
	if err != nil {
		log.Fatalf("Migration failed for key %s: insert failed: %s", err)
	}
	
	stmt.Close()
}

func GetBool(key string) bool {
	raw := GetString(key)
	return raw == "true"
}

func GetString(key string) string {
	var value string

	stmt := prepare("select Value from config where Key = ?")
	defer stmt.Close()

	row := stmt.QueryRow(key)
	err := row.Scan(&value)

	if err != nil {
		log.Fatalf("Unable to get configuration value with key %s: %s", key, err)
	}

	return value
}

func Set(key string, value interface{}) {
	stmt := prepare("delete from config where Key = ?")
	defer stmt.Close()

	_, err := stmt.Exec(key)
	if err != nil {
		log.Fatalf("Unable to set value for %s: %s", key, err)
	}

	insert := prepare("insert into config values (?, ?)")
	defer insert.Close()

	_, insErr := stmt.Exec(key, value)
	if insErr != nil {
		log.Fatalf("Unable to set value for %s: %s", key, err)
	}
}

func IsUser(username string) bool {
	var count int

	stmt := prepare("select count(*) from auth where Username = ?")
	defer stmt.Close()

	// will never error
	stmt.QueryRow(username).Scan(&count)

	return count != 0
}

func GetUser(raw string) structs.User {
	var username string
	var password string
	var twofactor string

	stmt := prepare("select * from auth where Username = ?")
	defer stmt.Close()

	err := stmt.QueryRow(raw).Scan(&username, &password, &twofactor)
	if err != nil {
		log.Fatalf("Unable to get user with name %s: %s", raw, err)
	}

	return structs.User {
		Username: username,
		Password: password,
		TwoFactor: twofactor,
	}
}

func GetUsers() map[string]structs.User {
	var users = make(map[string]structs.User)

	stmt := prepare("select * from auth")
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatalf("Unable to list users: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		var password string
		var twofactor string

		if err := rows.Scan(&username, &password, &twofactor); err != nil {
			log.Fatalf("Unable to list user: %s", err)
		}

		users[username] = structs.User {
			Username: username,
			Password: password,
			TwoFactor: twofactor,
		}
	}

	return users
}

// TODO: remove tx param after migration done
func CreateUser(username string, hash string, tx *sql.Tx) {
	stmt := prepare("insert into auth values (?, ?, '')")
	if tx != nil {
		stmt = tx.Stmt(stmt)
	}
	defer stmt.Close()
	
	stmt.Exec(username, hash)
}