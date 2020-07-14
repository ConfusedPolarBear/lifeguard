// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"log"
	"syscall"

	"github.com/ConfusedPolarBear/lifeguard/pkg/api"
	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
	"github.com/ConfusedPolarBear/lifeguard/pkg/crypto"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	resetFlag  := flag.String("r", "", "Username to reset password for")
	createFlag := flag.String("c", "", "Username to create")
	flag.Parse()
	reset := *resetFlag
	create := *createFlag

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.Normalize();

	log.Println("Starting Lifeguard")
	log.Printf("Git commit: %s%s", config.Commit, config.Modified)
	log.Printf("Build time: %s", config.BuildTime)
	log.Printf("Go version: %s", config.GoVersion)

	config.Load()

	hash := ""
	if reset != "" || create != "" {
		fmt.Print("Enter new password: ")

		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Unable to get password: %s", err)
		}
		fmt.Println()

		hash = crypto.HashPassword(string(bytePassword))
	}

	if reset != "" {
		config.SetPassword(reset, hash)
		log.Printf("Successfully reset password for %s", reset)
		return

	} else if create != "" {
		config.CreateUser(create, hash, nil)
		log.Printf("Successfully created account for %s", create)
		return
	}

	api.Setup()
}
