#!/usr/bin/env bash

# shellcheck disable=SC2034 # Variables sourced in other scripts.

# ------------
# NOTE: All functions that return lists should use newlines.
# bash functions can't return arrays, and spaces are tricky, so newline
# separators are the preferred pattern.
# To transform a string of newline-separated items to an array, use x::util::read-array:
# x::util::read-array FOO < <(x::golang::dups a b c a)
#
# ALWAYS remember to quote your subshells. Not doing so will break in
# bash 4.3, and potentially cause other issues.
# ------------

# The golang package that we are building.
X_GO_PACKAGE=github.com/fleezesd/xnightwatch

# Returns a sorted newline-separated list containing only duplicated items.
x::golang::dups() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort | uniq -d
}

# Returns a sorted newline-separated list with duplicated items removed.
x::golang::dedup() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort -u
}

# Asks golang what it thinks the host platform is. The go tool chain does some
# slightly different things when the target platform matches the host platform.
x::golang::host_platform() {
  echo "$(go env GOHOSTOS)_$(go env GOHOSTARCH)"
}

# Takes the platform name ($1) and sets the appropriate golang env variables
# for that platform.
x::golang::set_platform_envs() {
  [[ -n ${1-} ]] || {
    x::log::error_exit "!!! Internal error. No platform set in x::golang::set_platform_envs"
  }

  export GOOS=${platform%_*}
  export GOARCH=${platform##*_}

  # Do not set CC when building natively on a platform, only if cross-compiling
  if [[ $(x::golang::host_platform) != "$platform" ]]; then
    # Dynamic CGO linking for other server architectures than host architecture goes here
    # If you want to include support for more server platforms than these, add arch-specific gcc names here
    case "${platform}" in
      "linux_amd64")
        export CGO_ENABLED=1
        export CC=${KUBE_LINUX_AMD64_CC:-x86_64-linux-gnu-gcc}
        ;;
      "linux_arm")
        export CGO_ENABLED=1
        export CC=${KUBE_LINUX_ARM_CC:-arm-linux-gnueabihf-gcc}
        ;;
      "linux_arm64")
        export CGO_ENABLED=1
        export CC=${KUBE_LINUX_ARM64_CC:-aarch64-linux-gnu-gcc}
        ;;
      "linux_ppc64le")
        export CGO_ENABLED=1
        export CC=${KUBE_LINUX_PPC64LE_CC:-powerpc64le-linux-gnu-gcc}
        ;;
      "linux_s390x")
        export CGO_ENABLED=1
        export CC=${KUBE_LINUX_S390X_CC:-s390x-linux-gnu-gcc}
        ;;
    esac
  fi

  # if CC is defined for platform then always enable it
  ccenv=$(echo "$platform" | awk -F/ '{print "X_" toupper($1) "_" toupper($2) "_CC"}')
  if [ -n "${!ccenv-}" ]; then
    export CGO_ENABLED=1
    export CC="${!ccenv}"
  fi
}

# Ensure the go tool exists and is a viable version.
# Inputs:
#   env-var GO_VERSION is the desired go version to use, downloading it if needed (defaults to content of .go-version)
#   env-var FORCE_HOST_GO set to a non-empty value uses the go version in the $PATH and skips ensuring $GO_VERSION is used
x::golang::verify_go_version() {
  # default GO_VERSION to content of .go-version
  GO_VERSION="${GO_VERSION:-"$(cat "${X_ROOT}/.go-version")"}"
  if [ "${GOTOOLCHAIN:-auto}" != 'auto' ]; then
    # no-op, just respect GOTOOLCHAIN
    :
  elif [ -n "${FORCE_HOST_GO:-}" ]; then
    # ensure existing host version is used, like before GOTOOLCHAIN existed
    export GOTOOLCHAIN='local'
  else
    # otherwise, we want to ensure the go version matches GO_VERSION
    GOTOOLCHAIN="go${GO_VERSION}"
    export GOTOOLCHAIN
    # if go is either not installed or too old to respect GOTOOLCHAIN then use gimme
    if ! (command -v go >/dev/null && [ "$(go version | cut -d' ' -f3)" = "${GOTOOLCHAIN}" ]); then
      export GIMME_ENV_PREFIX=${GIMME_ENV_PREFIX:-"${LOCAL_OUTPUT_ROOT}/.gimme/envs"}
      export GIMME_VERSION_PREFIX=${GIMME_VERSION_PREFIX:-"${LOCAL_OUTPUT_ROOT}/.gimme/versions"}
      # eval because the output of this is shell to set PATH etc.
      eval "$("${X_ROOT}/third_party/gimme/gimme" "${GO_VERSION}")"
    fi
  fi

  if [[ -z "$(command -v go)" ]]; then
    x::log::usage_from_stdin <<EOF
Can't find 'go' in PATH, please fix and retry.
See http://golang.org/doc/install for installation instructions.
EOF
    return 2
  fi

  local go_version
  IFS=" " read -ra go_version <<< "$(go version)"
  local minimum_go_version
  minimum_go_version=go1.20
  if [[ "${minimum_go_version}" != $(echo -e "${minimum_go_version}\n${go_version[2]}" | sort -s -t. -k 1,1 -k 2,2n -k 3,3n | head -n1) && "${go_version[2]}" != "devel" ]]; then
    x::log::usage_from_stdin <<EOF
Detected go version: ${go_version[*]}.
OneX requires ${minimum_go_version} or greater.
Please install ${minimum_go_version} or later.
EOF
    return 2
  fi
}

# x::golang::setup_env will check that the `go` commands is available in
# ${PATH}. It will also check that the Go version is good enough for the
# Node build.
#
# Outputs:
#   env-var GOBIN is unset (we want binaries in a predictable place)
#   env-var GO15VENDOREXPERIMENT=1
#   env-var GO111MODULE=on
#   current directory is within GOPATH
x::golang::setup_env() {
  x::golang::verify_go_version

  # Set GOROOT so binaries that parse code can work properly.
  GOROOT=$(go env GOROOT)
  export GOROOT

  # Unset GOBIN in case it already exists in the current session.
  unset GOBIN

  # This seems to matter to some tools
  export GO15VENDOREXPERIMENT=1

  # Open go module feature
  export GO111MODULE=on

  # This is for sanity.  Without it, user umasks leak through into release
  # artifacts.
  umask 0022
}

# Args:
#  $1: platform (e.g. linux_amd64)
x::golang::build_binaries_for_platform() {
 # This is for sanity.  Without it, user umasks can leak through.
  umask 0022

  local platform=$1
  local command=$2

  V=2 x::log::info "Env for ${command}/${platform}: GOOS=${GOOS-} GOARCH=${GOARCH-} GOFLAGS=${GOFLAGS-} CGO_ENABLED=${CGO_ENABLED-} CC=${CC-}"
  V=2 x::log::info "Env for ${command}/${platform}: GOROOT=${GOROOT-}"
  V=3 x::log::info "Building binaries with GCFLAGS=${gogcflags} ASMFLAGS=${goasmflags} LDFLAGS=${goldflags}"

  local -a build_args
  build_args=(
    ${goflags:+"${goflags[@]}"}
    -gcflags="${gogcflags}"
    -asmflags="${goasmflags}"
    -ldflags="${goldflags}"
    -tags="${gotags:-}"
  )

  if [[ ${GOOS} == "windows" ]]; then
    ext="${bin}.exe"
  fi

  # Make sure the output directory exists
  out_dir="${LOCAL_OUTPUT_ROOT}/platforms/${GOOS}/${GOARCH}"
  mkdir -p ${out_dir}

  # Execute compilation command
  go build "${build_args[@]}" -o "${out_dir}/${command}${ext}" "${X_ROOT}/cmd/${command}"
  V=2 x::log::info "Output file is: ${out_dir}/${command}${ext}"
}


# Build binaries targets specified
#
# Args:
#  $1: command (e.g. onex-usercenter)
#  $2: platform (e.g. linux_amd64) (optional)
x::golang::build_binaries() {
  # Create a sub-shell so that we don't pollute the outer environment
  (
    # Check for `go` binary and set ${GOPATH}.
    x::golang::setup_env

    local command=$1
    local host_platform=$(x::golang::host_platform)
    local platform=${2:-${host_platform}}

    # These are "local" but are visible to and relied on by functions this
    # function calls.  They are effectively part of the calling API to
    # build_binaries_for_platform.
    local goflags goldflags goasmflags gogcflags gotags

    # This is $(pwd) because we use run-in-gopath to build.  Once that is
    # excised, this can become ${KUBE_ROOT}.
    local trimroot # two lines to appease shellcheck SC2155
    trimroot=$(pwd)

    goasmflags="all=-trimpath=${trimroot}"

    gogcflags="all=-trimpath=${trimroot} ${GOGCFLAGS:-}"
    if [[ "${DBG:-}" == 1 ]]; then
        # Debugging - disable optimizations and inlining and trimPath
        gogcflags="${GOGCFLAGS:-} all=-N -l"
        goasmflags=""
    fi

    goldflags="all=$(x::version::ldflags) ${GOLDFLAGS:-}"
    #goldflags="all=$(x::version::ldflags) ${GOLDFLAGS:-}"
    if [[ "${DBG:-}" != 1 ]]; then
        # Not debugging - disable symbols and DWARF.
        goldflags="${goldflags} -s -w"
    fi

    # Extract tags if any specified in GOFLAGS
    gotags="selinux,notest,$(echo "${GOFLAGS:-}" | sed -ne 's|.*-tags=\([^-]*\).*|\1|p')"

    x::golang::set_platform_envs "${platform}"
    x::log::status "${platform}: build started"
    x::golang::build_binaries_for_platform "${platform}" "${command}"
    x::log::status "${platform}: build finished"
  )
}
