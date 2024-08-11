#!/bin/bash -l

## Inputs

SPEC_FILE="$1"
# GITHUB_TOKEN="$2"
# UPLOAD_RELEASE_SCRIPT_URL='https://gist.github.com/stefanbuck/ce788fee19ab6eb0b4447a85fc99f447/raw/dbadd7d310ce8446de89c4ffdf1db0b400d0f6c3/upload-github-release-asset.sh'
#
# readonly GITHUB_TOKEN UPLOAD_RELEASE_SCRIPT_URL

## Build RPM

# Create RPM build tree
rpmdev-setuptree

# Move SPEC file to RPM tree
mv -v "$SPEC_FILE" ~/rpmbuild/SPECS/

# Find package metadata
SPEC_FILE="$HOME/rpmbuild/SPECS/$SPEC_FILE"
PKG_VERSION="$(rpmspec -q --qf "%{VERSION}\n" "$SPEC_FILE" | sort -u | head -n 1)"
PKG_NAME="$(rpmspec -q --qf '%{NAME}\n' "$SPEC_FILE" | grep -v -P '\-(devel|debug).*$' | sort -u | head -n 1)"

# Find Git metadata
GIT_RELEASE_TAG="$(basename "$GITHUB_REF")"
GIT_REPO_NAME="$(basename "$GITHUB_REPOSITORY")"

# Ensure package and git metadata match
if [[ ! v"$PKG_VERSION" == "$GIT_RELEASE_TAG" ]]; then
  echo "ERROR: Mismatch in package version and git release tag. PKG_VER: $PKG_VERSION; GIT_RELEASE_TAG: $GIT_RELEASE_TAG"
  exit 1
fi
if [[ ! "$PKG_NAME" == "$GIT_REPO_NAME" ]]; then
  echo "ERROR: Mismatch in package name and git repository name. PKG_NAME: $PKG_NAME; GIT_REPO_NAME: $GIT_REPO_NAME"
  exit 1
fi

readonly SPEC_FILE PKG_VERSION PKG_NAME GIT_RELEASE_TAG GIT_REPO_NAME

# Install build requirements
dnf --refresh -y builddep "$SPEC_FILE"

# Create archive
git archive --output=~/rpmbuild/SOURCES/"$PKG_NAME"-"$PKG_VERSION".tar.gz --prefix="$PKG_NAME"-"$PKG_VERSION"/ "$GIT_RELEASE_TAG"

# Build package
rpmbuild -ba "$SPEC_FILE"

# Define output
echo "packages_dir=~/rpmbuild/RPMS/" >>"$GITHUB_OUTPUT"

#
# ## Upload asset to Github Release
#
# # Download "upload-github-release-asset.sh"
# curl -X GET -L -o upload-github-release-asset.sh "$UPLOAD_RELEASE_SCRIPT_URL"
# bash ./upload-github-release-asset.sh \
#   github_api_token="$GITHUB_TOKEN" \
#   owner="$GITHUB_REPOSITORY_OWNER" \
#   repo="$GIT_REPO_NAME" \
#   tag="$GIT_RELEASE_TAG" \
#   filename=~/rpmbuild/RPMS/x86_64/"$PKG_NAME"-"$PKG_VERSION".el9.x86_64.rpm
