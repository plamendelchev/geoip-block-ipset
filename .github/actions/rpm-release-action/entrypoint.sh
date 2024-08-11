#!/bin/bash -l

## Inputs
SPEC_FILE="$1"

## Build RPM

# Create RPM build tree
RPMBUILD_DIR="$(realpath ./rpmbuild/)"
mkdir -v -p "$RPMBUILD_DIR"/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}

# Copy SPEC file to RPM tree
cp -v "$SPEC_FILE" "$RPMBUILD_DIR"/SPECS/

# Find package metadata
SPEC_FILE="$RPMBUILD_DIR/SPECS/$(basename "$SPEC_FILE")"
PKG_VERSION="$(rpmspec -q --qf "%{VERSION}\n" "$SPEC_FILE" | sort -u | head -n 1)"
PKG_NAME="$(rpmspec -q --qf '%{NAME}\n' "$SPEC_FILE" | grep -v -P '\-(devel|debug).*$' | sort -u | head -n 1)"

# Find Git metadata
GIT_RELEASE_TAG="$(basename "$GITHUB_REF")"
GIT_REPO_NAME="$(basename "$GITHUB_REPOSITORY")"

# Print all vars
echo "RPMBUILD_DIR: $RPMBUILD_DIR"
echo "SPEC_FILE: $SPEC_FILE"
echo "PKG_VERSION: $PKG_VERSION"
echo "PKG_NAME: $PKG_NAME"
echo "GIT_RELEASE_TAG: $GIT_RELEASE_TAG"
echo "GIT_REPO_NAME: $GIT_REPO_NAME"

# Ensure package and git metadata match
if [[ ! v"$PKG_VERSION" == "$GIT_RELEASE_TAG" ]]; then
  echo "ERROR: Mismatch in package version and git release tag. PKG_VER: $PKG_VERSION; GIT_RELEASE_TAG: $GIT_RELEASE_TAG"
  exit 1
fi
if [[ ! "$PKG_NAME" == "$GIT_REPO_NAME" ]]; then
  echo "ERROR: Mismatch in package name and git repository name. PKG_NAME: $PKG_NAME; GIT_REPO_NAME: $GIT_REPO_NAME"
  exit 1
fi

readonly RPMBUILD_DIR SPEC_FILE PKG_VERSION PKG_NAME GIT_RELEASE_TAG GIT_REPO_NAME

# Install build requirements
dnf --refresh -y builddep "$SPEC_FILE" || exit 1

# Fix for "detected dubious ownership in repository"
git config --global --add safe.directory /github/workspace

# Create archive
git archive --output="$RPMBUILD_DIR"/SOURCES/"$PKG_NAME"-"$PKG_VERSION".tar.gz --prefix="$PKG_NAME"-"$PKG_VERSION"/ "$GIT_RELEASE_TAG" || exit 1

# Build package
rpmbuild --define "_topdir $RPMBUILD_DIR" -ba "$SPEC_FILE" || exit 1

# Define output
echo "packages_dir=./rpmbuild/RPMS" >>"$GITHUB_OUTPUT"
