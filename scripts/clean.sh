#!/usr/bin/env bash

	echo "Clean up ..."
	rm -rf /tmp/enroll_user /tmp/msp kidner
	docker rm -f -v `docker ps -a --no-trunc | grep "kidner" | cut -d ' ' -f 1` 2>/dev/null || true
	docker rmi `docker images --no-trunc | grep "kidner" | cut -d ' ' -f 1` 2>/dev/null || true
	echo "Clean up done"