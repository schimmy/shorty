#!/bin/bash
set -ex

export REDIS_URL=$REDIS_PORT_6379_TCP_ADDR:$REDIS_PORT_6379_TCP_PORT
export PG_USER=postgres
export PG_HOST=localhost
export PG_DATABASE=shortener

createdb -h "$PG_HOST" -U "$PG_USER" "$PG_DATABASE"
psql -U $PG_USER -d $PG_DATABASE -h "$PG_HOST" < ./pg_schema.sql

make test
