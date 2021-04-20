#!/bin/bash

# unzip ./data/test_accounts_*.zip

PGPASSWORD=postgres psql -h localhost -d postgres -U postgres -p 5431 \
    -c "\i migrations/schema.sql"

go build ./cmd/dataloader
./dataloader --conn "postgres://postgres:postgres@localhost:5431/postgres?sslmode=disable" \
    ./data/data/accounts_1.json