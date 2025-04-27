#!/usr/bin/env bash

# This file is not intended to be run automatically. It is meant to be run
# immediately before exporting docs. We do not want to check these documents in
# by default.

set -o errexit
set -o nounset
set -o pipefail

X_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${X_ROOT}/scripts/lib/init.sh"

x::golang::setup_env
x::util::ensure-temp-dir

BINS=(
  gen-docs
  gen-x-docs
  gen-man
  gen-yaml
)
make build -C "${X_ROOT}" BINS="${BINS[*]}"

# Run all known doc generators (today gendocs and genman for nodectl)
# $1 is the directory to put those generated documents
function generate_docs() {
  local dest="$1"

  # Find binary
  gendocs=$(x::util::find-binary "gen-docs")
  genxdocs=$(x::util::find-binary "gen-x-docs")
  genman=$(x::util::find-binary "gen-man")
  genyaml=$(x::util::find-binary "gen-yaml")

  mkdir -p "${dest}/docs/guide/en-US/cmd/xctl"
  "${gendocs}" "${dest}/docs/guide/en-US/cmd/xctl/"

  mkdir -p "${dest}/docs/guide/en-US/cmd"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/" "x-fakeserver"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/" "x-usercenter"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/" "x-apiserver"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/" "x-gateway"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/" "x-nightwatch"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/" "x-pump"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/" "x-toyblc"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/" "x-controller-manager"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/" "x-minerset-controller"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/" "x-miner-controller"
  "${genxdocs}" "${dest}/docs/guide/en-US/cmd/xctl" "xctl"

  mkdir -p "${dest}/docs/man/man1/"
  "${genman}" "${dest}/docs/man/man1/" "x-fakeserver"
  "${genman}" "${dest}/docs/man/man1/" "x-usercenter"
  "${genman}" "${dest}/docs/man/man1/" "x-apiserver"
  "${genman}" "${dest}/docs/man/man1/" "x-gateway"
  "${genman}" "${dest}/docs/man/man1/" "x-nightwatch"
  "${genman}" "${dest}/docs/man/man1/" "x-pump"
  "${genman}" "${dest}/docs/man/man1/" "x-toyblc"
  "${genman}" "${dest}/docs/man/man1/" "x-controller-manager"
  "${genman}" "${dest}/docs/man/man1/" "x-minerset-controller"
  "${genman}" "${dest}/docs/man/man1/" "x-miner-controller"
  "${genman}" "${dest}/docs/man/man1/" "xctl"

  mkdir -p "${dest}/docs/guide/en-US/yaml/xctl/"
  "${genyaml}" "${dest}/docs/guide/en-US/yaml/xctl/"

  # create the list of generated files
  pushd "${dest}" > /dev/null || return 1
  touch docs/.generated_docs
  find . -type f | cut -sd / -f 2- | LC_ALL=C sort > docs/.generated_docs
  popd > /dev/null || return 1
}

# Removes previously generated docs-- we don't want to check them in. $X_ROOT
# must be set.
function remove_generated_docs() {
  if [ -e "${X_ROOT}/docs/.generated_docs" ]; then
    # remove all of the old docs; we don't want to check them in.
    while read -r file; do
      rm "${X_ROOT}/${file}" 2>/dev/null || true
    done <"${X_ROOT}/docs/.generated_docs"
    # The docs/.generated_docs file lists itself, so we don't need to explicitly
    # delete it.
  fi
}

# generate into X_TMP
generate_docs "${X_TEMP}"

# remove all of the existing docs in x_ROOT
remove_generated_docs

# Copy fresh docs into the repo.
# the shopt is so that we get docs/.generated_docs from the glob.
shopt -s dotglob
cp -af "${X_TEMP}"/* "${X_ROOT}"
shopt -u dotglob