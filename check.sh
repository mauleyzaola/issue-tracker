#!/bin/bash
echo "

trying to download all go dependencies...
"

if ! go get -v ./...; then
    echo "[FAIL] Failed to get dependencies."
    exit 1
fi

go get bitbucket.org/liamstask/goose/cmd/goose
go get github.com/stretchr/testify/assert

echo "

upgrading test database objects to latest version...
"

cd dbmigrations/pg

if ! goose -env="test" up; then
    echo "[FAIL] Failed to upgrade the test database"
    exit 1
fi

cd ../../


echo "

testing all the packages...
"
cd test

if ! go test -v ./...; then
    echo "[FAIL] One or more unit tests failed."
    exit 1
fi

cd ..
