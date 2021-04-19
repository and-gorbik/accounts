#!/bin/bash

# unzip ./data/test_accounts_*.zip

go build ./cmd/dataloader
./dataloader ./data/data/accounts_1.json --conn "user=postgres password=postgres host=localhost port=5431 database=postgres sslmode=disable"