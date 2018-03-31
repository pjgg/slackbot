#!/bin/bash -x

source ci/build.env.sh


$DOCKER build --no-cache -t $DOCKER_IMAGE . 

