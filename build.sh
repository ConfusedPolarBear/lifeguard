#!/bin/bash
set -euo pipefail

# base config path
export CONFIG="github.com/ConfusedPolarBear/lifeguard/pkg/config"

# build time
export TIME="$(date)"

# go version
export VERSION="$(go version)"

# returns " (modified)" if local changes have been made to the git repository, "" otherwise
modified="$(git diff --no-ext-diff --quiet || echo " (modified)")"
commit="$(git rev-list -1 HEAD)"
commit="${commit:0:7}"

# git commit this build was created from
export COMMIT="$commit$modified"

echo "Step 1: Building server"
go build -ldflags "-X '$CONFIG.Commit=$COMMIT' -X '$CONFIG.BuildTime=$TIME' -X '$CONFIG.GoVersion=$VERSION'"

echo
echo "Step 2: Building web UI"
npm run dev

echo
echo "Successfully built lifeguard"
