// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"sync"

	"github.com/ConfusedPolarBear/lifeguard/pkg/config"

	"golang.org/x/crypto/bcrypt"
)

var known sync.Map

func GenerateHMAC(plaintext string) string {
	data := []byte(plaintext)
	secret := config.GetString("security.session_key")
	h := hmac.New(sha256.New, []byte(secret))

	h.Write(data)

	hmac := hex.EncodeToString(h.Sum(nil))
	known.Store(hmac, plaintext)

	return hmac
}

func LookupHMAC(hmac string) string {
	// Load() is probably vulnerable to a timing attack but irrelevant since the HMAC values aren't secret.
	// HMACs are only used to prevent command injection.
	if value, ok := known.Load(hmac); ok {
		return value.(string)
	}

	return ""
}

func HashPassword(plaintext string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		log.Fatalf("Unable to hash password: %s", err)
	}

	return string(hash)
}

func GetRandom(n int) string {
	b := generateRandomBytes(n)
	return hex.EncodeToString(b)
}

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Unable to generate random bytes: %s", err)
	}

	return b
}
