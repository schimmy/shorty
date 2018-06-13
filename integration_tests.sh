#!/bin/bash
set -ex

export PG_USER=postgres
export PG_HOST=localhost
export PG_DATABASE=shortener

createdb -h "$PG_HOST" -U "$PG_USER" "$PG_DATABASE"
psql -U $PG_USER -d $PG_DATABASE -h "$PG_HOST" < ./pg_schema.sql

make test
