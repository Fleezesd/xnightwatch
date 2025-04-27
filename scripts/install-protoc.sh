#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

X_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${X_ROOT}/scripts/lib/init.sh"
source "${X_ROOT}/scripts/lib/protoc.sh"

x::protoc::install