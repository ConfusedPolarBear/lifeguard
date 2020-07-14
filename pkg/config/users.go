// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"database/sql"
	"log"

	"github.com/ConfusedPolarBear/lifeguard/pkg/structs"

	_ "github.com/mattn/go-sqlite3"
)

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

func SetPassword(username string, hash string) {
	stmt := prepare("update auth set Password = ? where Username = ?")
	defer stmt.Close()
	
	stmt.Exec(hash, username)
}