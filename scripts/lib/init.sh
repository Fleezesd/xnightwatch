#!/usr/bin/env bash

set -o errexit
set +o nounset
set -o pipefail

# Short-circuit if init.sh has already been sourced
[[ $(type -t x::init::loaded) == function ]] && return 0

# Unset CDPATH so that path interpolation can work correctly
# https://github.com/minerrnetes/minerrnetes/issues/52255
unset CDPATH

# Default use go modules
export GO111MODULE=on

# The root of the build/dist directory
X_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

X_OUTPUT_SUBPATH="${X_OUTPUT_SUBPATH:-_output}"
X_OUTPUT="${X_ROOT}/${X_OUTPUT_SUBPATH}"

source "${X_ROOT}/scripts/lib/util.sh"
source "${X_ROOT}/scripts/lib/logging.sh"
source "${X_ROOT}/scripts/lib/color.sh"

x::log::install_errexit
x::util::ensure-bash-version

source "${X_ROOT}/scripts/lib/version.sh"
source "${X_ROOT}/scripts/lib/golang.sh"

# list of all available group versions. This should be used when generated code
# or when starting an API server that you want to have everything.
# most preferred version for a group should appear first
# UPDATEME: New group need to update here.
X_AVAILABLE_GROUP_VERSIONS="${X_AVAILABLE_GROUP_VERSIONS:-\
apps/v1beta1 \
}"

# This emulates "readlink -f" which is not available on MacOS X.
# Test:
# T=/tmp/$$.$RANDOM
# mkdir $T
# touch $T/file
# mkdir $T/dir
# ln -s $T/file $T/linkfile
# ln -s $T/dir $T/linkdir
# function testone() {
#   X=$(readlink -f $1 2>&1)
#   Y=$(x::readlinkdashf $1 2>&1)
#   if [ "$X" != "$Y" ]; then
#     echo readlinkdashf $1: expected "$X", got "$Y"
#   fi
# }
# testone /
# testone /tmp
# testone $T
# testone $T/file
# testone $T/dir
# testone $T/linkfile
# testone $T/linkdir
# testone $T/nonexistant
# testone $T/linkdir/file
# testone $T/linkdir/dir
# testone $T/linkdir/linkfile
# testone $T/linkdir/linkdir
function x::readlinkdashf {
  # run in a subshell for simpler 'cd'
  (
    if [[ -d "${1}" ]]; then # This also catch symlinks to dirs.
      cd "${1}"
      pwd -P
    else
      cd "$(dirname "${1}")"
      local f
      f=$(basename "${1}")
      if [[ -L "${f}" ]]; then
        readlink "${f}"
      else
        echo "$(pwd -P)/${f}"
      fi
    fi
  )
}

# This emulates "realpath" which is not available on MacOS X
# Test:
# T=/tmp/$$.$RANDOM
# mkdir $T
# touch $T/file
# mkdir $T/dir
# ln -s $T/file $T/linkfile
# ln -s $T/dir $T/linkdir
# function testone() {
#   X=$(realpath $1 2>&1)
#   Y=$(x::realpath $1 2>&1)
#   if [ "$X" != "$Y" ]; then
#     echo realpath $1: expected "$X", got "$Y"
#   fi
# }
# testone /
# testone /tmp
# testone $T
# testone $T/file
# testone $T/dir
# testone $T/linkfile
# testone $T/linkdir
# testone $T/nonexistant
# testone $T/linkdir/file
# testone $T/linkdir/dir
# testone $T/linkdir/linkfile
# testone $T/linkdir/linkdir
x::realpath() {
  if [[ ! -e "${1}" ]]; then
    echo "${1}: No such file or directory" >&2
    return 1
  fi
  x::readlinkdashf "${1}"
}

# Marker function to indicate init.sh has been fully sourced
x::init::loaded() {
  return 0
}
