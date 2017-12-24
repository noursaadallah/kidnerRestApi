#!/usr/bin/env bash

	echo "Build ..."
	govendor sync
	go build
	echo "Build done"