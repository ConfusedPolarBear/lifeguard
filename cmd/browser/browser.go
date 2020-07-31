// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

/* browser is the backend program that powers Lifeguard's file browser. Since Lifeguard does not run as root and ships
 * with it's own authentication layer, it needs to be able to access any file in a ZFS dataset (or snapshot) for the
 * file browser to function correctly. To ensure that a malicious user on the system (or someone acting througn an RCE)
 * cannot use this binary to read any file on the system, the browser executable has multiple security layers in place:
 *   - It first checks that the config file is owned (and only writable) by root and that it permits browsing files.
 *   - The path specified with the "-f" flag is verified against the config file.
 *   - If the both previous checks pass, browser verifies that the parent process is the same Lifeguard executable that
*        the browser was compiled with.
*/

package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/viper"
)

var DisableSecurity string
var SHA256 string
var SHA512 string

var IsProduction bool = (DisableSecurity != "yes")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flagPath     := flag.String("f", "", "Path to browse to")
	flag.Parse()
	path := *flagPath

	log.Printf("Initializing Lifeguard file browser")

	// Print baked in hashes of Lifeguard or security warning if in dev mode
	if IsProduction {
		log.Printf("Accept: %s/%s", SHA256[:10], SHA512[:10])
	} else {
		log.Println()
		log.Printf("This binary was built in development mode, which disables core security features")
		log.Printf("that prevent arbitrary file read vulnerabilities. Do *not* disable them in production!")
		log.Println()
	}

	// Run safety checks
	if os.Geteuid() != 0 {
		ReportCheckFail("EUID is not root")
	}
	AssertConfigPermissions()
	AssertHashes()

	// Parse the configuration
	viper.SetConfigName("browser")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Unable to load config: %s", err)
	}

	// Verify that file browsing is enabled and allowed for the given path
	if !viper.GetBool("browser.enabled") {
		log.Fatalf("File browser is disabled")
	}

	all := viper.GetString("browser.allowed")
	ok := false
	for _, allowed := range strings.Split(all, ",") {
		if strings.HasPrefix(path, allowed) {
			ok = true
			break
		}
	}
	if !ok {
		log.Fatalf("Path %s is not allowed", path)
	}

	log.Printf("Lifeguard file browser initialized")

	// Open and stat the path
	file, openErr := os.Open(path)
	if openErr != nil {
		log.Fatalf("Error: unable to open %s: %s", path, openErr)
	}
	
	info, statErr := file.Stat()
	if statErr != nil {
		log.Fatalf("Error: unable to stat %s: %s", path, statErr)
	}

	// List the contents of the directory or read the file
	if info.IsDir() {
		contents, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatalf("Error: unable to list directory contents of %s: %s", path, err)
		}

		fmt.Printf("fold")
		for _, entry := range contents {
			typeChar := "f"
			name := base64.StdEncoding.EncodeToString([]byte(entry.Name()))

			if entry.IsDir() {
				typeChar = "d"
			}

			// Output format: "type encoded_name bytes"
			fmt.Printf("%s %s %d\n", typeChar, name, entry.Size())
		}
	} else {
		fmt.Printf("file")
		if _, err := io.Copy(os.Stdout, file); err != nil {
			log.Fatalf("Error: Unable to copy %s: %s", path, err)
		}
	}
}

func AssertConfigPermissions() {
	f, openErr := os.Open("config/browser.ini")
	if openErr != nil {
		log.Fatalf("Error: unable to open config: %s", openErr)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		log.Fatalf("Error: unable to stat config: %s", err)
	}

	uid := info.Sys().(*syscall.Stat_t).Uid
	perms := fmt.Sprintf("%#o", info.Mode().Perm())

	if uid != 0 {
		ReportCheckFail("config file must be owned by root")
	}

	if perms != "0600" {
		ReportCheckFail("config file must have 0600 permissions")
	}
}

func AssertHashes() {
	ppid := strconv.FormatInt(int64(os.Getppid()), 10)

	// Read the parent's executable to hash and verify it
	f, err := os.Open("/proc/" + string(ppid) + "/exe")
	if err != nil {
		log.Fatalf("Error: unable to open parent process executable: %s", err)
	}
	defer f.Close()

	sha256 := sha256.New()
	sha512 := sha512.New()

	if _, err := io.Copy(sha256, f); err != nil {
		log.Fatalf("Error: unable to copy file to sha256 instance: %s", err)
	}

	f.Seek(0, 0)

	if _, err := io.Copy(sha512, f); err != nil {
		log.Fatalf("Error: unable to copy file to sha512 instance: %s", err)
	}

	calculated256 := fmt.Sprintf("%x", sha256.Sum(nil))
	calculated512 := fmt.Sprintf("%x", sha512.Sum(nil))

	if SHA256 != calculated256 || SHA512 != calculated512 {
		ReportCheckFail("checksum of caller " + ppid + " incorrect")

		if SHA256 != "" {
			log.Println()
			log.Printf("calc SHA256:  %s", calculated256)
			log.Printf("known SHA256: %s", SHA256)
			log.Printf("calc SHA512:  %s", calculated512)
			log.Printf("known SHA512: %s", SHA512)
			log.Println()
		}
	}
}

func ReportCheckFail(msg string) {
	level := "Error"

	if IsProduction {
		log.Fatalf("%s: %s", level, msg)
	} else {
		level = "Warning"
		log.Printf("%s: %s (would be fatal in production)", level, msg)
	}
}
