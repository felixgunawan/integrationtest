#!/bin/bash

docker-compose up -d
go build -o integrationtest ./test
./integrationtest