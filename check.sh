#!/bin/bash
echo "

trying to download all go dependencies...
"

if ! go get -v ./...; then
    echo "[FAIL] Failed to get dependencies."
    exit 1
fi

go get github.com/rubenv/sql-migrate/...
go get github.com/stretchr/testify/assert

echo "

upgrading test database objects to latest version...
"

psql -c 'drop database if exists tracker_test;'
psql -c 'create database tracker_test;'

if [ ! -f migrations/dbconfig.yml ]; then
	cp migrations/dbconfig.yml.sample migrations/dbconfig.yml
fi	

if [ ! -f test/config.json ]; then
	cp test/config.json.sample test/config.json
fi	

if [ ! -f server/config.json ]; then
	cp server/config.json.sample server/config.json
fi	

echo "

testing all the packages...
"
cd test

#cannot run all tests at once, so we loop instead on each directory looking for test files
for dir in ./*
do
    dir=${dir%*/}
    if [ -d "$dir" ]; then
	cd "$dir"
	if ! go test -v; then
	   echo "[FAIL] One or more unit tests failed."
	   exit 1
	fi 
	cd ..
    fi
done

cd ..
