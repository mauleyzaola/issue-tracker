#!/bin/bash
sudo service nginx stop
sudo stop issue-tracker

psql -c 'drop database tracker;'
psql -c 'create database tracker;'

pg_restore -d tracker tracker.sample.backup


cd $GOPATH/src/github.com/mauleyzaola/issue-tracker

git pull
cd server/
go build


sudo service nginx start
sudo start issue-tracker
