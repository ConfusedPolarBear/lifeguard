// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pquerna/otp/totp"
)

func InitializeTOTP(username string) (string, string) {
	secret, _ := totp.Generate(totp.GenerateOpts {
		Issuer: "Lifeguard",
		AccountName: username,
	})

	image, err := secret.Image(256, 256)
	if err != nil {
		log.Printf("Unable to create TOTP qr code: %s", err)
		return "", ""
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, image)

	return secret.Secret(), base64.StdEncoding.EncodeToString(buf.Bytes())
}

func SaveTOTP(username string, secret string, code string) bool {
	if !totp.Validate(code, secret) {
		return false
	}

	stmt := prepare("update auth set TwoFactorProvider = 'totp', TwoFactorData = ? where Username = ?")
	defer stmt.Close()
	
	stmt.Exec(secret, username)

	return true
}

func VerifyTOTP(username string, code string) bool {
	var secret string

	stmt := prepare("select TwoFactorData from auth where Username = ?")
	defer stmt.Close()

	err := stmt.QueryRow(username).Scan(&secret)
	if err != nil {
		log.Fatalf("Unable to get TOTP secret for user %s: %s", username, err)
	}

	if !totp.Validate(code, secret) {
		log.Printf("%s failed 2FA challenge: Invalid TOTP code", username)
		return false
	}

	log.Printf("%s authenticated successfully", username)
	return true
}