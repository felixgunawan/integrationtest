#!/bin/bash

docker-compose -f docker-development/docker-compose.yml up -d
go build -o example 
go build -o integrationtest ./test
./wait-for-it.sh localhost:5432 -- ./example &
./wait-for-it.sh localhost:55001 -- ./integrationtest
kill $(lsof -t -i:55001)