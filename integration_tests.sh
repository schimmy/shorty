#!/bin/bash

export REDIS_URL=localhost:7777

createdb -h localhost -U postgres drone
export PG_HOST=localhost
export PG_DATABASE=drone
psql -h "$PG_HOST" < ./pg_schema.sql


go test -v ./...

# kill all child processes to clean up
pkill -P $$
