.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@sh ./scripts/build.sh

##### ENV
env-up:
	@sh ./scripts/env_up.sh

env-down:
	@sh ./scripts/env_down.sh

##### RUN
run:
	@echo "Start app ..."
	@./kidnerRestApi

##### CLEAN
clean: env-down
	@sh ./scripts/clean.sh

#### init the prerequisites
setup-preq:
	@sh ./scripts/setupPRQ.sh

#### init the setup : hyperledger fabric & fabric-ca & fabric-sdk-go
setup-hf:
	@sh ./scripts/setupHF.sh

#### help
help:
#	@chmod a+x ./scripts/*
	
	@echo "Usage : "
	@echo
	@echo "setup-preq ..... install latest version of Docker and docker-compose and install go v1.8.3"
	@echo "                 (you need to logout/login after this step for the changes to take effect)"
	@echo
	@echo "setup-hf ....... install hyperledger fabric, fabric-ca, and fabric-sdk-go (v1.0.0-rc1) and install external libs"
	@echo
	@echo "all ............ clean the environment, build the app, start environment and run"
	@echo
	@echo "dev ............ build and run the app"
	@echo
	@echo "build .......... build the app"
	@echo
	@echo "run ............ run the app"
	@echo
	@echo "clean .......... clean-up temporary files and docker images"
	@echo
	@echo "env-up ......... start the necessary docker images"
	@echo
	@echo "env-down ....... stop the running docker images" 
	@echo