#!/usr/bin/env bash

set -e

mkdir -p release
cd release

echo "Created /release directory."

build() {
    os=$1
    arch=$2

    if [ ${os} = darwin ]; then
        buildos="macos"
    else
        buildos=${os}
    fi

    if [ ${arch} = amd64 ]; then
        buildarch="64bit"
    else
        buildarch="32bit"
    fi

    binary=oro
    release="$binary-$buildos-$buildarch"

    if [ ${os} = windows ]; then
        binary="$binary.exe"
    fi

    env GOOS=${os} GOARCH=${arch} go build -v -o ${binary} ../oro.go

    if [ ${os} = linux ]; then
        tar czf "$release.tar.gz" "$binary"
    else
        zip "$release.zip" "$binary"
    fi

    rm -f ${binary}

    echo "Created release for '$buildos' on '$buildarch'."
}

# MacOS
build darwin amd64
build darwin 386

# Linux
build linux amd64
build linux 386

# Windows
build windows amd64
build windows 386