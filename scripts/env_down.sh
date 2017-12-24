#!/usr/bin/env bash
	
    echo "Stop environnement ..."
	cd fixtures && docker-compose down
	echo "Environnement down"