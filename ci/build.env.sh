#!/bin/bash



##############################################################
##                      CONFIGURE THIS                      ##
##############################################################

export SRC_PATH="go/src/github.com/pjgg/slackbot"

########################### END ##############################

export GIT_TAG=$(git describe --abbrev=0 --tags 2> /dev/null)
export GIT_BRANCH=$(basename $(git branch -r --contains HEAD))
export GIT_COMMIT=$(git rev-parse --short HEAD)
export GIT_COMMIT_SHORT=$(git rev-parse HEAD)

if [ -z $GIT_TAG ]
then 
    export BASE_VERSION="0.0.1";
else 
    export BASE_VERSION=$GIT_TAG
fi


if [ "$GIT_BRANCH" = "master" ]
then 
    export VERSION=$BASE_VERSION
else 
    export VERSION=$BASE_VERSION"-SNAPSHOT"
fi

# Check if can run docker without sudo
docker ps > /dev/null
if [ $? -eq 0 ]; then
    export DOCKER="docker"
    export GCLOUD="gcloud"
else
    export DOCKER="sudo docker"
    export GCLOUD="sudo gcloud"
fi

export DOCKER_IMAGE="pjgg/slackbot:$VERSION"
