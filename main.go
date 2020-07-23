// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"log"
	"syscall"
	"os"

	"github.com/ConfusedPolarBear/lifeguard/pkg/api"
	"github.com/ConfusedPolarBear/lifeguard/pkg/config"
	"github.com/ConfusedPolarBear/lifeguard/pkg/crypto"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	devFlag    := flag.Bool("d", false, "Enable debug mode")
	resetFlag  := flag.String("r", "", "Username to reset password for")
	createFlag := flag.String("c", "", "Username to create")
	tfaFlag    := flag.String("t", "", "Username to remove 2FA for")
	flag.Parse()
	config.DevMode = *devFlag
	reset := *resetFlag
	create := *createFlag
	tfa := *tfaFlag

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.Normalize();

	log.Println("Starting Lifeguard")
	log.Printf("Git commit: %s%s", config.Commit, config.Modified)
	log.Printf("Build time: %s", config.BuildTime)
	log.Printf("Go version: %s", config.GoVersion)
	if config.DevMode {
		log.Printf("Warning: Debug mode has been enabled. Some security features will be disabled.")
	}

	config.Load()

	// Command line operations should only be available to root
	prompt := reset != "" || create != ""
	if (prompt || tfa != "") && os.Geteuid() != 0 {
		log.Fatalf("CLI is only available to root")
	}

	hash := ""
	if prompt {
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

	} else if tfa != "" {
		config.DisableTwoFactor(tfa)
		log.Printf("Successfully disabled two factor for %s", tfa)
		return
	}

	api.Setup()
}
