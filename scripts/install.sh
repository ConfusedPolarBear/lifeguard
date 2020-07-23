#!/bin/bash
# Copyright Matt Montgomery (licensed AGPLv3)

set -eu

log() {
	echo "$(date) $*"
}

if [[ $# -ne 2 ]]; then
	log "Usage: $0 lifeguardUser poolName"
fi

user="$1"		# User that lifeguard will run as
pool="$2"		# Pool name to permit lifeguard to access

log "Installing Lifeguard"

if [[ "$EUID" -ne "0" ]]; then
	log "Error: this script must be run as root"
	exit 1
fi

# =========== sudo policy ===========
log "Installing new sudo policy as 99-lifeguard"
sed "s/%LG_USER%/$user/" assets/99-lifeguard > /etc/sudoers.d/99-lifeguard

# =========== zpool permissions ===========
log "Allowing user $user to perform the following actions on pool $pool:"
log "    View the differences between a dataset and it's snapshot (diff)"
log "    Load and unload encryption keys (load-key)"
log "    Mount and unmount datasets (mount)"
log "    Create snapshots of datasets (snapshot)"
zfs allow -d -u "$user" diff,load-key,mount,snapshot "$pool"

# =========== browser permissions ===========
log "Setting browser permissions and ownership"
# The browser binary needs to be SUID root and r-xr-xr-x
chown root browser
chmod 4555 browser

# The browser config needs to be only readable/writable by root and rw-------
chown root config/browser.ini
chmod 0600 config/browser.ini

# =========== done ===========
log "Installation successful"