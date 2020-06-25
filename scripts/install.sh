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

log Installing new sudo policy as 99-lifeguard
cp assets/99-lifeguard /etc/sudoers.d/
