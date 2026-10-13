#!/bin/sh
set -eu

if [ "${MIGRATION_RUN:-true}" = "true" ]; then
  attempts=0
  until ./go-template-lepsios migrate up; do
    attempts=$((attempts + 1))
    if [ "$attempts" -ge "30" ]; then
      echo "database migrations failed after ${attempts} attempts" >&2
      exit 1
    fi
    echo "database is not ready; retrying migrations (${attempts}/30)" >&2
    sleep 2
  done
fi

exec ./go-template-lepsios server
