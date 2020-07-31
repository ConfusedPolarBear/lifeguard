// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func prepare(raw string) *sql.Stmt {
	stmt, err := db.Prepare(raw)
	if err != nil {
		log.Fatalf("Unable to prepare statement %s: %s", raw, err)
	}

	return stmt
}

func GetTwoFactorProvider(username string) string {
	var provider string

	stmt := prepare("select TwoFactorProvider from auth where Username = ?")
	defer stmt.Close()

	if err := stmt.QueryRow(username).Scan(&provider); err != nil {
		log.Fatalf("Unable to get TOTP provider for user %s: %s", username, err)
	}

	return provider
}

func IsTwoFactorEnabled(username string) bool {
	return GetTwoFactorProvider(username) != ""
}

func DisableTwoFactor(username string) {
	stmt := prepare("update auth set TwoFactorProvider = '', TwoFactorData = '' where Username = ?")
	defer stmt.Close()
	
	stmt.Exec(username)	
}
