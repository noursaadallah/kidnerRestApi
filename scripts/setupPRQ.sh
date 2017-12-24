#!/usr/bin/env bash

	echo "Setting up the technologies required ..."

	# test if docker 17.x installed
	which docker >&2
	if [ $? -eq 0 ]
	then
		docker -v | grep "17" >&2
		if [ $? -eq 0 ]
		then
			echo "docker v17.x installed"
		else
			echo "installing docker ..."
			sudo apt-get update && sudo apt install docker.io
		fi
	else
		echo "installing docker ..."
		sudo apt-get update && sudo apt install docker.io
	fi 
		
		# test if docker-compose 1.14.x installed
	which docker-compose >&2
	if [ $? -eq 0 ]
	then
		docker-compose -v | grep "1.14" >&2
		if [ $? -eq 0 ]
		then
			echo "docker-compose v1.14.x installed"
		else
			echo "installing docker-compose ..."
			sudo apt-get update && sudo apt install docker-compose
		fi
	else
		echo "installing docker-compose ..."
		sudo apt-get update && sudo apt install docker-compose
	fi 
		
		echo "configure docker user to current user"
		sudo groupadd docker
		sudo gpasswd -a ${USER} docker
		sudo service docker restart
	

	#test if go 1.8.x is installed 
	which go >&2
	if [ $? -eq 0 ]
	then 
		go version | grep "go1.8" >&2
		if [ $? -eq 0 ]
			then
				echo "go 1.8.x installed"
			else
				echo "install Go 1.8.3 and set the GOPATH to $HOME/go"
				wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
				sudo tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz
				rm go1.8.3.linux-amd64.tar.gz
				echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
				echo 'export GOPATH=$HOME/go' | tee -a $HOME/.bashrc
				echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' | tee -a $HOME/.bashrc
				mkdir -p $HOME/go/{src,pkg,bin}
			fi	
	else
		echo "install Go 1.8.3 and set the GOPATH to $HOME/go"
		wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
		sudo tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz
		rm go1.8.3.linux-amd64.tar.gz
		echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
		echo 'export GOPATH=$HOME/go' | tee -a $HOME/.bashrc
		echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' | tee -a $HOME/.bashrc
		mkdir -p $HOME/go/{src,pkg,bin}
	fi
    
	echo "please login and logout for the changes to take effect"