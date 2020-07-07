#!/bin/bash
# Copyright Matt Montgomery (licensed AGPLv3)

set -eu

log() {
	echo "$(date) $@"
}

log Installing Lifeguard

if [[ "$EUID" -ne "0" ]]; then
	log This script must be run as root.
	exit 1
fi

# =========== sudo policy ===========
log Installing new sudo policy as 99-lifeguard
cp assets/99-lifeguard /etc/sudoers.d/

# =========== browser permissions ===========
log Setting browser permissions and ownership
# The browser binary needs to be SUID root and r-xr-xr-x
chown root browser
chmod 4555 browser

# The browser config needs to be only readable/writable by root and rw-------
chown root config/browser.ini
chmod 0600 config/browser.ini

# =========== done ===========
log Installation successful