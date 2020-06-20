#!/bin/bash
set -euo pipefail

export TIME="$(date)"
export VERSION="$(go version)"
export CONFIG="github.com/ConfusedPolarBear/lifeguard/pkg/config"

commit="$(git rev-list -1 HEAD)"
commit="${commit:0:7}"
export COMMIT="$commit"

echo "Step 1: Building server"
go build -ldflags "-X '$CONFIG.Commit=$COMMIT' -X '$CONFIG.BuildTime=$TIME' -X '$CONFIG.GoVersion=$VERSION'"

echo
echo "Step 2: Building web UI"
npm run dev

echo
echo "Successfully built lifeguard"
