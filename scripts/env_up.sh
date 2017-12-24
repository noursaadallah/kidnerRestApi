#!/usr/bin/env bash

    echo "Start environment ..."
	cd fixtures && docker-compose up --force-recreate -d
	echo "Sleep 10 seconds in order to let the environment setup correctly"
	sleep 10
	echo "Environment up"