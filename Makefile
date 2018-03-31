VERSION = $(shell  if [ $TAG = "" ]; then echo "0.0.1"; else echo "$TAG"; fi)
SRC_PATH=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SLACK_TOKEN=YourSlackBotToken

version:
	echo $(VERSION)

clean: 
	rm -rf\
	 $(SRC_PATH)/dist\
	 $(SRC_PATH)/debug\
	 $(SRC_PATH)/godepgraph.png\
	 $(SRC_PATH)/*/cover.out\
	 $(SRC_PATH)/*/cover.html\
	 $(SRC_PATH)/pkg\
	 $(SRC_PATH)/src\

build:
	ci/build-go.sh

build-docker:
	ci/build-docker.sh

run-local: all
	docker run -d -e CONFIG_PATH=config -e SLACK_TOKEN=$(SLACK_TOKEN) pjgg/slackbot:0.0.1-SNAPSHOT /app slack-listener

all: clean build build-docker
