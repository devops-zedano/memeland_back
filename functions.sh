#!/bin/bash
set -e

# Set DEBUG=true to enable debugging of the script.
if [ "$DEBUG" = "true" ]; then set -x; fi

build_func() {
    local func_name; func_name="$1"
    local workspace; workspace="$2"

    mkdir -p "${wprkspace}/bin"

    echo "Building function ${func_name}...."

    cd "$func_name"
    go build -o "$workspace/bin"
    cd $workspace
}

build_all() {
    cwd=$(pwd)

    for func_name in $(ls -d */ | grep func); do
        build_func "$func_name" "$cwd"
    done

    cd $cwd
}

deploy_all() {
    echo 'Deploying the whole application...'
    serverless deploy
    rm -rf bin/*
}

remove() {
    echo 'Destroying the whole application...'
    serverless remove
}

deploy_func() {
    local func_name; func_name="$1"
    echo "Deploying function $func_name..."
    serverless deploy function --function "$func_name" 
}



