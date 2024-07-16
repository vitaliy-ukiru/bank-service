#!/usr/bin/env bash

set -e

APP=""


function migrate () {
  if [[ ! -z "${RUN_MIGRATIONS}" ]]; then
    ./bin/migrator
  fi
}

migrate
./bin/web-server $*
