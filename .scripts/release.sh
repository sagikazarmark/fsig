#!/bin/bash

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "The new version parameter is required"

    exit 1
fi

# Increase version in README installation instructions
sed -i -e "s/ENV FSIG_VERSION .*/ENV FSIG_VERSION $VERSION/g" README.md
rm -f README.md-e

sed -i -e "s/^## \[Unreleased\]$/## [Unreleased]\\"$'\n'"\\"$'\n'"\\"$'\n'"## [${VERSION}] - $(date +%Y-%m-%d)/g" CHANGELOG.md
rm -f CHANGELOG.md-e

sed -i -e "s|^\[Unreleased\]: \(.*\)HEAD$|[Unreleased]: https://github.com/sagikazarmark/fsig/compare/v${VERSION}...HEAD\\"$'\n'"[${VERSION}]: \1v${VERSION}|g" CHANGELOG.md
rm -f CHANGELOG.md-e
