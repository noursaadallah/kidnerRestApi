#!/usr/bin/env bash

	echo "setup hyperledger fabric & fabric-ca & fabric-sdk-go & external libs"
	echo ""
	
    echo "setup fabric to version v1.0.0-rc1"
	echo ""

	if [ -d "$GOPATH/src/github.com/hyperledger/fabric" ]
	then
		cd $GOPATH/src/github.com/hyperledger/fabric && git checkout v1.0.0-rc1
	fi
	if [ -d "$GOPATH/src/github.com/hyperledger/fabric-ca" ]
	then
		cd $GOPATH/src/github.com/hyperledger/fabric-ca && git checkout v1.0.0-rc1
	fi
	if [ -d "$GOPATH/src/github.com/hyperledger/fabric-sdk-go" ]
	then
		cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go && git checkout 85fa3101eb4694d464003c3a900672d632f17833
		echo "installing fabric-sdk-go dependencies"
		echo ""
		cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go && make
	fi

	
	if [ ! -d "$GOPATH/src/github.com/hyperledger" ]
	then
		mkdir -p $GOPATH/src/github.com/hyperledger
		cd $GOPATH/src/github.com/hyperledger && git clone https://github.com/hyperledger/fabric.git && cd fabric && git checkout v1.0.0-rc1
		echo ""
		echo "setup fabric-ca to version v1.0.0-rc1"
		echo ""
		cd $GOPATH/src/github.com/hyperledger && git clone https://github.com/hyperledger/fabric-ca.git && cd fabric-ca && git checkout v1.0.0-rc1

		echo ""	
		echo "setup fabric-sdk-go to commit 85fa3101eb4694d464003c3a900672d632f17833"
		echo ""
		cd $GOPATH/src/github.com/hyperledger && git clone https://github.com/hyperledger/fabric-sdk-go.git && cd fabric-sdk-go && git checkout 85fa3101eb4694d464003c3a900672d632f17833
		
		echo ""
		echo "installing fabric-sdk-go dependencies"
		echo ""
		cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go && make
		
		echo ""
		echo "fabric & fabric-ca & fabric-sdk-go setup done"
		echo ""
		echo "installing fabric-sdk-go client packages"
		echo ""
		sudo apt install libltdl-dev
		go get github.com/hyperledger/fabric-sdk-go/pkg/fabric-client
		go get github.com/hyperledger/fabric-sdk-go/pkg/fabric-ca-client

		echo ""
		echo "setup external libs"
		echo ""
		go get -u github.com/spf13/viper
		go get -u github.com/kardianos/govendor
		# make govendor executable just in case $PATH doesn't include $GOPATH/bin
		sudo cp $GOPATH/bin/govendor /usr/local/bin/ 
		cd $GOPATH/src/github.com/noursaadallah/kidner && govendor init && govendor add +external
	fi

	echo ""
	echo "setup external libs"
	echo ""
	go get -u github.com/spf13/viper
	go get -u github.com/kardianos/govendor
	sudo cp $GOPATH/bin/govendor /usr/local/bin/
	cd $GOPATH/src/github.com/noursaadallah/kidner && govendor init && govendor add +external