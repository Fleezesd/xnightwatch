#!/usr/bin/env bash

# This file creates release artifacts (tar files, container images) that are
# ready to distribute to install or distribute to end users.

###############################################################################
# Most of the ::release:: namespace functions have been moved to
# github.com/onex/release.  Have a look in that repo and specifically in
# lib/releaselib.sh for ::release::-related functionality.
###############################################################################

# Tencent cos configuration
readonly BUCKET="superproj-1254073058"
readonly REGION="ap-beijing"
readonly COS_RELEASE_DIR="onex-release"
readonly COSTOOL="coscmd"

# This is where the final release artifacts are created locally
readonly RELEASE_STAGE="${LOCAL_OUTPUT_ROOT}/release-stage"
readonly RELEASE_TARS="${LOCAL_OUTPUT_ROOT}/release-tars"
readonly RELEASE_IMAGES="${LOCAL_OUTPUT_ROOT}/release-images"

# onex github account info
readonly X_GITHUB_ORG=superproj
readonly X_GITHUB_REPO=onex

readonly ARTIFACT=onex.tar.gz
readonly CHECKSUM=${ARTIFACT}.sha1sum

X_BUILD_CONFORMANCE=${X_BUILD_CONFORMANCE:-y}
X_BUILD_PULL_LATEST_IMAGES=${X_BUILD_PULL_LATEST_IMAGES:-y}

# Validate a ci version
#
# Globals:
#   None
# Arguments:
#   version
# Returns:
#   If version is a valid ci version
# Sets:                    (e.g. for '1.2.3-alpha.4.56+abcdef12345678')
#   VERSION_MAJOR          (e.g. '1')
#   VERSION_MINOR          (e.g. '2')
#   VERSION_PATCH          (e.g. '3')
#   VERSION_PRERELEASE     (e.g. 'alpha')
#   VERSION_PRERELEASE_REV (e.g. '4')
#   VERSION_BUILD_INFO     (e.g. '.56+abcdef12345678')
#   VERSION_COMMITS        (e.g. '56')
function x::release::parse_and_validate_ci_version() {
  # Accept things like "v1.2.3-alpha.4.56+abcdef12345678" or "v1.2.3-beta.4"
  local -r version_regex="^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)-([a-zA-Z0-9]+)\\.(0|[1-9][0-9]*)(\\.(0|[1-9][0-9]*)\\+[0-9a-f]{7,40})?$"
  local -r version="${1-}"
  [[ "${version}" =~ ${version_regex} ]] || {
    x::log::error "Invalid ci version: '${version}', must match regex ${version_regex}"
    return 1
  }

  # The VERSION variables are used when this file is sourced, hence
  # the shellcheck SC2034 'appears unused' warning is to be ignored.

  # shellcheck disable=SC2034
  VERSION_MAJOR="${BASH_REMATCH[1]}"
  # shellcheck disable=SC2034
  VERSION_MINOR="${BASH_REMATCH[2]}"
  # shellcheck disable=SC2034
  VERSION_PATCH="${BASH_REMATCH[3]}"
  # shellcheck disable=SC2034
  VERSION_PRERELEASE="${BASH_REMATCH[4]}"
  # shellcheck disable=SC2034
  VERSION_PRERELEASE_REV="${BASH_REMATCH[5]}"
  # shellcheck disable=SC2034
  VERSION_BUILD_INFO="${BASH_REMATCH[6]}"
  # shellcheck disable=SC2034
  VERSION_COMMITS="${BASH_REMATCH[7]}"
}

# ---------------------------------------------------------------------------
# Build final release artifacts
function x::release::clean_cruft() {
  # Clean out cruft
  find "${RELEASE_STAGE}" -name '*~' -exec rm {} \;
  find "${RELEASE_STAGE}" -name '#*#' -exec rm {} \;
  find "${RELEASE_STAGE}" -name '.DS*' -exec rm {} \;
}

function x::release::package_tarballs() {
  # Clean out any old releases
  rm -rf "${RELEASE_STAGE}" "${RELEASE_TARS}" "${RELEASE_IMAGES}"
  mkdir -p "${RELEASE_TARS}"
  x::release::package_src_tarball &
  x::release::package_client_tarballs &
  x::release::package_X_manifests_tarball &
  x::release::package_server_tarballs &
  x::util::wait-for-jobs || { x::log::error "previous tarball phase failed"; return 1; }

  x::release::package_final_tarball & # _final depends on some of the previous phases
  x::util::wait-for-jobs || { x::log::error "previous tarball phase failed"; return 1; }
}

function x::release::updload_tarballs() {
  x::log::info "upload ${RELEASE_TARS}/* to cos bucket ${BUCKET}."
  for file in $(ls ${RELEASE_TARS}/*)
  do
    if [ "${COSTOOL}" == "coscli" ];then
      coscli cp "${file}" "cos://${BUCKET}/${COS_RELEASE_DIR}/${X_GIT_VERSION}/"
      coscli cp "${file}" "cos://${BUCKET}/${COS_RELEASE_DIR}/latest/"
    else
      coscmd upload  "${file}" "${COS_RELEASE_DIR}/${X_GIT_VERSION}/"
      coscmd upload  "${file}" "${COS_RELEASE_DIR}/latest/"
    fi
  done
}

# Package the source code we built, for compliance/licensing/audit/yadda.
function x::release::package_src_tarball() {
  local -r src_tarball="${RELEASE_TARS}/onex-src.tar.gz"
  x::log::status "Building tarball: src"
  if [[ "${X_GIT_TREE_STATE-}" = 'clean' ]]; then
    git archive -o "${src_tarball}" HEAD
  else
    find "${X_ROOT}" -mindepth 1 -maxdepth 1 \
      ! \( \
      \( -path "${X_ROOT}"/_\* -o \
      -path "${X_ROOT}"/.git\* -o \
      -path "${X_ROOT}"/.gitignore\* -o \
      -path "${X_ROOT}"/.gsemver.yaml\* -o \
      -path "${X_ROOT}"/.config\* -o \
      -path "${X_ROOT}"/.chglog\* -o \
      -path "${X_ROOT}"/.gitlint -o \
      -path "${X_ROOT}"/.golangci.yaml -o \
      -path "${X_ROOT}"/.goreleaser.yml -o \
      -path "${X_ROOT}"/.note.md -o \
      -path "${X_ROOT}"/.todo.md \
      \) -prune \
      \) -print0 \
      | "${TAR}" czf "${src_tarball}" --transform "s|${X_ROOT#/*}|onex|" --null -T -
  fi
}

# Package up all of the server binaries
function x::release::package_server_tarballs() {
  # Find all of the built client binaries
  local long_platforms=("${LOCAL_OUTPUT_BINPATH}"/*/*)
  if [[ -n ${X_BUILD_PLATFORMS-} ]]; then
    read -ra long_platforms <<< "${X_BUILD_PLATFORMS}"
  fi

  for platform_long in "${long_platforms[@]}"; do
    local platform
    local platform_tag
    platform=${platform_long##${LOCAL_OUTPUT_BINPATH}/} # Strip LOCAL_OUTPUT_BINPATH
    platform_tag=${platform/\//-} # Replace a "/" for a "-"
    x::log::status "Starting tarball: server $platform_tag"

    (
    local release_stage="${RELEASE_STAGE}/server/${platform_tag}/onex"
    rm -rf "${release_stage}"
    mkdir -p "${release_stage}/server/bin"

    local server_bins=("${X_SERVER_BINARIES[@]}")

      # This fancy expression will expand to prepend a path
      # (${LOCAL_OUTPUT_BINPATH}/${platform}/) to every item in the
      # server_bins array.
      cp "${server_bins[@]/#/${LOCAL_OUTPUT_BINPATH}/${platform}/}" \
        "${release_stage}/server/bin/"

      x::release::clean_cruft

      local package_name="${RELEASE_TARS}/onex-server-${platform_tag}.tar.gz"
      x::release::create_tarball "${package_name}" "${release_stage}/.."
      ) &
    done

    x::log::status "Waiting on tarballs"
    x::util::wait-for-jobs || { x::log::error "server tarball creation failed"; exit 1; }
  }

# Package up all of the cross compiled clients. Over time this should grow into
# a full SDK
function x::release::package_client_tarballs() {
  # Find all of the built client binaries
  local long_platforms=("${LOCAL_OUTPUT_BINPATH}"/*/*)
  if [[ -n ${X_BUILD_PLATFORMS-} ]]; then
    read -ra long_platforms <<< "${X_BUILD_PLATFORMS}"
  fi

  for platform_long in "${long_platforms[@]}"; do
    local platform
    local platform_tag
    platform=${platform_long##${LOCAL_OUTPUT_BINPATH}/} # Strip LOCAL_OUTPUT_BINPATH
    platform_tag=${platform/\//-} # Replace a "/" for a "-"
    x::log::status "Starting tarball: client $platform_tag"

    (
    local release_stage="${RELEASE_STAGE}/client/${platform_tag}/onex"
    rm -rf "${release_stage}"
    mkdir -p "${release_stage}/client/bin"

    local client_bins=("${X_CLIENT_BINARIES[@]}")

      # This fancy expression will expand to prepend a path
      # (${LOCAL_OUTPUT_BINPATH}/${platform}/) to every item in the
      # client_bins array.
      cp "${client_bins[@]/#/${LOCAL_OUTPUT_BINPATH}/${platform}/}" \
        "${release_stage}/client/bin/"

      x::release::clean_cruft

      local package_name="${RELEASE_TARS}/onex-client-${platform_tag}.tar.gz"
      x::release::create_tarball "${package_name}" "${release_stage}/.."
    ) &
  done

  x::log::status "Waiting on tarballs"
  x::util::wait-for-jobs || { x::log::error "client tarball creation failed"; exit 1; }
}

# Package up all of the server binaries in docker images
function x::release::build_server_images() {
  # Clean out any old images
  rm -rf "${RELEASE_IMAGES}"
  local platform
  for platform in "${X_SERVER_PLATFORMS[@]}"; do
    local platform_tag
    local arch
    platform_tag=${platform/\//-} # Replace a "/" for a "-"
    arch=$(basename "${platform}")
    x::log::status "Building images: $platform_tag"

    local release_stage
    release_stage="${RELEASE_STAGE}/server/${platform_tag}/onex"
    rm -rf "${release_stage}"
    mkdir -p "${release_stage}/server/bin"

    # This fancy expression will expand to prepend a path
    # (${LOCAL_OUTPUT_BINPATH}/${platform}/) to every item in the
    # X_SERVER_IMAGE_BINARIES array.
    cp "${X_SERVER_IMAGE_BINARIES[@]/#/${LOCAL_OUTPUT_BINPATH}/${platform}/}" \
      "${release_stage}/server/bin/"

    x::release::create_docker_images_for_server "${release_stage}/server/bin" "${arch}"
  done
}

function x::release::md5() {
  if which md5 >/dev/null 2>&1; then
    md5 -q "$1"
  else
    md5sum "$1" | awk '{ print $1 }'
  fi
}

function x::release::sha1() {
  if which sha1sum >/dev/null 2>&1; then
    sha1sum "$1" | awk '{ print $1 }'
  else
    shasum -a1 "$1" | awk '{ print $1 }'
  fi
}

function x::release::build_conformance_image() {
  local -r arch="$1"
  local -r registry="$2"
  local -r version="$3"
  local -r save_dir="${4-}"
  x::log::status "Building conformance image for arch: ${arch}"
  ARCH="${arch}" REGISTRY="${registry}" VERSION="${version}" \
    make -C cluster/images/conformance/ build >/dev/null

  local conformance_tag
  conformance_tag="${registry}/conformance-${arch}:${version}"
  if [[ -n "${save_dir}" ]]; then
    "${DOCKER[@]}" save "${conformance_tag}" > "${save_dir}/conformance-${arch}.tar"
  fi
  x::log::status "Deleting conformance image ${conformance_tag}"
  "${DOCKER[@]}" rmi "${conformance_tag}" &>/dev/null || true
}

# This builds all the release docker images (One docker image per binary)
# Args:
#  $1 - binary_dir, the directory to save the tared images to.
#  $2 - arch, architecture for which we are building docker images.
function x::release::create_docker_images_for_server() {
  # Create a sub-shell so that we don't pollute the outer environment
  (
    local binary_dir
    local arch
    local binaries
    local images_dir
    binary_dir="$1"
    arch="$2"
    binaries=$(x::build::get_docker_wrapped_binaries "${arch}")
    images_dir="${RELEASE_IMAGES}/${arch}"
    mkdir -p "${images_dir}"

    # k8s.gcr.io is the constant tag in the docker archives, this is also the default for config scripts in GKE.
    # We can use X_DOCKER_REGISTRY to include and extra registry in the docker archive.
    # If we use X_DOCKER_REGISTRY="k8s.gcr.io", then the extra tag (same) is ignored, see release_docker_image_tag below.
    local -r docker_registry="k8s.gcr.io"
    # Docker tags cannot contain '+'
    local docker_tag="${X_GIT_VERSION/+/_}"
    if [[ -z "${docker_tag}" ]]; then
      x::log::error "git version information missing; cannot create Docker tag"
      return 1
    fi

    # provide `--pull` argument to `docker build` if `X_BUILD_PULL_LATEST_IMAGES`
    # is set to y or Y; otherwise try to build the image without forcefully
    # pulling the latest base image.
    local docker_build_opts
    docker_build_opts=
    if [[ "${X_BUILD_PULL_LATEST_IMAGES}" =~ [yY] ]]; then
        docker_build_opts='--pull'
    fi

    for wrappable in $binaries; do

      local binary_name=${wrappable%%,*}
      local base_image=${wrappable##*,}
      local binary_file_path="${binary_dir}/${binary_name}"
      local docker_build_path="${binary_file_path}.dockerbuild"
      local docker_file_path="${docker_build_path}/Dockerfile"
      local docker_image_tag="${docker_registry}/${binary_name}-${arch}:${docker_tag}"

      x::log::status "Starting docker build for image: ${binary_name}-${arch}"
      (
        rm -rf "${docker_build_path}"
        mkdir -p "${docker_build_path}"
        ln "${binary_file_path}" "${docker_build_path}/${binary_name}"
        ln "${X_ROOT}/build/nsswitch.conf" "${docker_build_path}/nsswitch.conf"
        chmod 0644 "${docker_build_path}/nsswitch.conf"
        cat <<EOF > "${docker_file_path}"
FROM ${base_image}
COPY ${binary_name} /usr/local/bin/${binary_name}
EOF
        # ensure /etc/nsswitch.conf exists so go's resolver respects /etc/hosts
        if [[ "${base_image}" =~ busybox ]]; then
          echo "COPY nsswitch.conf /etc/" >> "${docker_file_path}"
        fi

        "${DOCKER[@]}" build ${docker_build_opts:+"${docker_build_opts}"} -q -t "${docker_image_tag}" "${docker_build_path}" >/dev/null
        # If we are building an official/alpha/beta release we want to keep
        # docker images and tag them appropriately.
        local -r release_docker_image_tag="${X_DOCKER_REGISTRY-$docker_registry}/${binary_name}-${arch}:${X_DOCKER_IMAGE_TAG-$docker_tag}"
        if [[ "${release_docker_image_tag}" != "${docker_image_tag}" ]]; then
          x::log::status "Tagging docker image ${docker_image_tag} as ${release_docker_image_tag}"
          "${DOCKER[@]}" rmi "${release_docker_image_tag}" 2>/dev/null || true
          "${DOCKER[@]}" tag "${docker_image_tag}" "${release_docker_image_tag}" 2>/dev/null
        fi
        "${DOCKER[@]}" save -o "${binary_file_path}.tar" "${docker_image_tag}" "${release_docker_image_tag}"
        echo "${docker_tag}" > "${binary_file_path}.docker_tag"
        rm -rf "${docker_build_path}"
        ln "${binary_file_path}.tar" "${images_dir}/"

        x::log::status "Deleting docker image ${docker_image_tag}"
        "${DOCKER[@]}" rmi "${docker_image_tag}" &>/dev/null || true
      ) &
    done

    if [[ "${X_BUILD_CONFORMANCE}" =~ [yY] ]]; then
      x::release::build_conformance_image "${arch}" "${docker_registry}" \
        "${docker_tag}" "${images_dir}" &
    fi

    x::util::wait-for-jobs || { x::log::error "previous Docker build failed"; return 1; }
    x::log::status "Docker builds done"
  )

}

# This will pack onex-system manifests files for distros such as COS.
function x::release::package_X_manifests_tarball() {
  x::log::status "Building tarball: manifests"

  local src_dir="${X_ROOT}/deployments"

  local release_stage="${RELEASE_STAGE}/manifests/onex"
  rm -rf "${release_stage}"

  local dst_dir="${release_stage}"
  mkdir -p "${dst_dir}"
  cp -r ${src_dir}/* "${dst_dir}"
  #cp "${src_dir}/onex-apiserver.yaml" "${dst_dir}"
  #cp "${src_dir}/onex-authz-server.yaml" "${dst_dir}"
  #cp "${src_dir}/onex-pump.yaml" "${dst_dir}"
  #cp "${src_dir}/onex-watcher.yaml" "${dst_dir}"
  #cp "${X_ROOT}/cluster/gce/gci/health-monitor.sh" "${dst_dir}/health-monitor.sh"

  x::release::clean_cruft

  local package_name="${RELEASE_TARS}/onex-manifests.tar.gz"
  x::release::create_tarball "${package_name}" "${release_stage}/.."
}

# This is all the platform-independent stuff you need to run/install onex.
# Arch-specific binaries will need to be downloaded separately (possibly by
# using the bundled cluster/get-onex-binaries.sh script).
# Included in this tarball:
#   - Cluster spin up/down scripts and configs for various cloud providers
#   - Tarballs for manifest configs that are ready to be uploaded
#   - Examples (which may or may not still work)
#   - The remnants of the docs/ directory
function x::release::package_final_tarball() {
  x::log::status "Building tarball: final"

  # This isn't a "full" tarball anymore, but the release lib still expects
  # artifacts under "full/onex/"
  local release_stage="${RELEASE_STAGE}/full/onex"
  rm -rf "${release_stage}"
  mkdir -p "${release_stage}"

  mkdir -p "${release_stage}/client"
  cat <<EOF > "${release_stage}/client/README"
Client binaries are no longer included in the onex final tarball.

Run release/get-onex-binaries.sh to download client and server binaries.
EOF

  # We want everything in /scripts.
  mkdir -p "${release_stage}/release"
  cp -R "${X_ROOT}/scripts/release" "${release_stage}/"
  cat <<EOF > "${release_stage}/release/get-onex-binaries.sh"
#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# This file download onex client and server binaries from tencent cos bucket.

os=linux arch=amd64 version=${X_GIT_VERSION} && wget https://${BUCKET}.cos.${REGION}.myqcloud.com/${COS_RELEASE_DIR}/\$version/{onex-client-\$os-\$arch.tar.gz,onex-server-\$os-\$arch.tar.gz}
EOF
  chmod +x ${release_stage}/release/get-onex-binaries.sh

  mkdir -p "${release_stage}/server"
  cp "${RELEASE_TARS}/onex-manifests.tar.gz" "${release_stage}/server/"
  cat <<EOF > "${release_stage}/server/README"
Server binary tarballs are no longer included in the onex final tarball.

Run release/get-onex-binaries.sh to download client and server binaries.
EOF

  # Include hack/lib as a dependency for the cluster/ scripts
  #mkdir -p "${release_stage}/hack"
  #cp -R "${X_ROOT}/hack/lib" "${release_stage}/hack/"

  cp -R ${X_ROOT}/{docs,configs,scripts,deployments,init,README.md,LICENSE} "${release_stage}/"

  echo "${X_GIT_VERSION}" > "${release_stage}/version"

  x::release::clean_cruft

  local package_name="${RELEASE_TARS}/${ARTIFACT}"
  x::release::create_tarball "${package_name}" "${release_stage}/.."
}

# Build a release tarball.  $1 is the output tar name.  $2 is the base directory
# of the files to be packaged.  This assumes that ${2}/onexis what is
# being packaged.
function x::release::create_tarball() {
  x::build::ensure_tar

  local tarfile=$1
  local stagingdir=$2

  "${TAR}" czf "${tarfile}" -C "${stagingdir}" onex --owner=0 --group=0
}

function x::release::install_github_release(){
  GO111MODULE=on go install github.com/github-release/github-release@latest
}

# Require the following tools:
# - github-release
# - gsemver
# - git-chglog
# - coscmd or coscli
function x::release::verify_prereqs(){
  if [ -z "$(which github-release 2>/dev/null)" ]; then
    x::log::info "'github-release' tool not installed, try to install it."

    if ! x::release::install_github_release; then
      x::log::error "failed to install 'github-release'"
      return 1
    fi
  fi

  if [ -z "$(which git-chglog 2>/dev/null)" ]; then
    x::log::info "'git-chglog' tool not installed, try to install it."

    if ! go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest &>/dev/null; then
      x::log::error "failed to install 'git-chglog'"
      return 1
    fi
  fi

  if [ -z "$(which gsemver 2>/dev/null)" ]; then
    x::log::info "'gsemver' tool not installed, try to install it."

    if ! go install github.com/arnaud-deprez/gsemver@latest &>/dev/null; then
      x::log::error "failed to install 'gsemver'"
      return 1
    fi
  fi


  if [ -z "$(which ${COSTOOL} 2>/dev/null)" ]; then
    x::log::info "${COSTOOL} tool not installed, try to install it."

    if ! make -C "${X_ROOT}" tools.install.${COSTOOL}; then
      x::log::error "failed to install ${COSTOOL}"
      return 1
    fi
  fi

  if [ -z "${TENCENT_SECRET_ID}" -o -z "${TENCENT_SECRET_KEY}" ];then
      x::log::error "can not find env: TENCENT_SECRET_ID and TENCENT_SECRET_KEY"
      return 1
  fi

  if [ "${COSTOOL}" == "coscli" ];then
    if [ ! -f "${HOME}/.cos.yaml" ];then
      cat << EOF > "${HOME}/.cos.yaml"
cos:
  base:
    secretid: ${TENCENT_SECRET_ID}
    secretkey: ${TENCENT_SECRET_KEY}
    sessiontoken: ""
  buckets:
  - name: ${BUCKET}
    alias: ${BUCKET}
    region: ${REGION}
EOF
    fi
  else
    if [ ! -f "${HOME}/.cos.conf" ];then
      cat << EOF > "${HOME}/.cos.conf"
[common]
secret_id = ${TENCENT_SECRET_ID}
secret_key = ${TENCENT_SECRET_KEY}
bucket = ${BUCKET}
region =${REGION}
max_thread = 5
part_size = 1
schema = https
EOF
    fi
  fi
}

# Create a github release with specified tarballs.
# NOTICE: Must export 'GITHUB_TOKEN' env in the shell, details:
# https://github.com/github-release/github-release
function x::release::github_release() {
  # create a github release
  x::log::info "create a new github release with tag ${X_GIT_VERSION}"
  github-release release \
    --user ${X_GITHUB_ORG} \
    --repo ${X_GITHUB_REPO} \
    --tag ${X_GIT_VERSION} \
    --description "" \
    --pre-release

  # update onex tarballs
  x::log::info "upload ${ARTIFACT} to release ${X_GIT_VERSION}"
  github-release upload \
    --user ${X_GITHUB_ORG} \
    --repo ${X_GITHUB_REPO} \
    --tag ${X_GIT_VERSION} \
    --name ${ARTIFACT} \
    --file ${RELEASE_TARS}/${ARTIFACT}

  x::log::info "upload onex-src.tar.gz to release ${X_GIT_VERSION}"
  github-release upload \
    --user ${X_GITHUB_ORG} \
    --repo ${X_GITHUB_REPO} \
    --tag ${X_GIT_VERSION} \
    --name "onex-src.tar.gz" \
    --file ${RELEASE_TARS}/onex-src.tar.gz
}

function x::release::generate_changelog() {
  x::log::info "generate CHANGELOG-${X_GIT_VERSION#v}.md and commit it"

  git-chglog ${X_GIT_VERSION} > ${X_ROOT}/CHANGELOG/CHANGELOG-${X_GIT_VERSION#v}.md

  set +o errexit
  git add ${X_ROOT}/CHANGELOG/CHANGELOG-${X_GIT_VERSION#v}.md
  git commit -a -m "docs(changelog): add CHANGELOG-${X_GIT_VERSION#v}.md"
  git push -f origin master # 最后将 CHANGELOG 也 push 上去
}

