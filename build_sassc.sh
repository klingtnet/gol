#!/usr/bin/env bash

TEMPDIR=$(mktemp --directory)
DEST=${1:-"./bin"}

if [ ! -d "$DEST" ]; then
    mkdir "$DEST"
fi

if [ $? -eq 0 ]; then
    curl -L "https://github.com/sass/libsass/archive/master.tar.gz" | tar -xvz -C "$TEMPDIR"
    curl -L "https://github.com/sass/sassc/archive/master.tar.gz" | tar -xvz -C "$TEMPDIR"
    pushd $TEMPDIR/sassc-master
    SASS_LIBSASS_PATH="$TEMPDIR/libsass-master" make build-static
    popd
    cp --verbose "$TEMPDIR/sassc-master/bin/sassc" "$DEST/sassc"
else
    echo "Could not create temporary directory using 'mktemp'!"
    exit 1
fi
