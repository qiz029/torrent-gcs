#!/bin/sh

echo "start to make dir"
mkdir /torrents-store

echo "start to map gsc and gce disk"
gcsfuse torrents-store /torrents-store

echo "resolve dependency"
dep ensure

echo "start the program"
go run main.go