## Install kidner :

    mkdir -p $GOPATH/src/github.com/noursaadallah
    cd $GOPATH/src/github.com/noursaadallah
    git clone https://github.com/noursaadallah/kidner.git
    cd kidner
    make help => display help
    make setup-preq => install prerequisites : docker & docker-compose & go v1.8.3
    make setup-hf => install fabric & fabric-ca & fabric-sdk-go
    make all => clean build env-up run

## Start kidner:

    cd $GOPATH/src/github.com/noursaadallah/kidner
    make env-up run


## ** Manual installation :

1. Prerequisites

    Ubuntu 16.04
    Go version 1.7.x or greater
    Docker version 1.12 or greater
    Docker-compose version 1.8 or greater

2. Hyperledger Fabric & Fabric-CA & Fabric-SDK-GO:

    2.1. Fabric v1.0.0-rc1:

        mkdir -p $GOPATH/src/github.com/hyperledger && \
        cd $GOPATH/src/github.com/hyperledger && \
        git clone https://github.com/hyperledger/fabric.git && \
        cd fabric && \
        git checkout v1.0.0-rc1

    2.2. Fabric-ca v1.0.0-rc1:

        cd $GOPATH/src/github.com/hyperledger && \
        git clone https://github.com/hyperledger/fabric-ca.git && \
        cd fabric-ca && \
        git checkout v1.0.0-rc1

    2.3. Fabric-SDK-GO :

        cd $GOPATH/src/github.com/hyperledger && \
        git clone https://github.com/hyperledger/fabric-sdk-go.git && \
        cd fabric-sdk-go && \
        git checkout 85fa3101eb4694d464003c3a900672d632f17833

        Install the client (you need sudo apt install libltdl-dev) :

            go get github.com/hyperledger/fabric-sdk-go/pkg/fabric-client && \
            go get github.com/hyperledger/fabric-sdk-go/pkg/fabric-ca-client

3. Prepare the environment :

    mkdir -p $GOPATH/src/github.com/noursaadallah && \
    cd $GOPATH/src/github.com/noursaadallah

    git clone https://github.com/noursaadallah/kidner.git

    make or make all 