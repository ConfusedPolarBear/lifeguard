// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
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
	if value, ok := known.Load(hmac); ok {
		return value.(string)
	}

	return ""
}
