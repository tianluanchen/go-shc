#!/usr/bin/env bash

go mod tidy && go mod download

build_with_os_arch() {
    export GOOS=$1
    export GOARCH=$2
    output=$3
    if [ -z "$output" ]; then
        output="$(basename $(go list -m))_#OS#_#ARCH#"
    fi
    echo "building for $(go env GOOS)/$(go env GOARCH)"
    suffix=""
    if [ "$GOOS" = "windows" ]; then
        suffix=".exe"
    fi
    output=$(echo "${output}${suffix}" | sed "s/#OS#/${GOOS}/g" | sed "s/#ARCH#/${GOARCH}/g")
    go build -ldflags "-s -w " \
        -gcflags="all=-trimpath=${PWD}" \
        -asmflags="all=-trimpath=${PWD}" \
        -o "${output}"
}
name="./bin/$(basename $(go list -m))_#OS#_#ARCH#"

if [ -d "bin" ]; then
    rm bin/*
fi

build_with_os_arch linux amd64 "$name" &
build_with_os_arch linux arm64 "$name" &
build_with_os_arch windows amd64 "$name" &
build_with_os_arch windows arm64 "$name" &
build_with_os_arch darwin amd64 "$name" &
build_with_os_arch darwin arm64 "$name" &
wait
echo "fininshed"
