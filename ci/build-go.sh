#!/bin/sh


go get -u -v github.com/kardianos/govendor



pwd
govendor sync

for GOOS in linux; do
  for GOARCH in amd64; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    go build -ldflags "-w -s" -o bin/app-$GOOS-$GOARCH
  done
done

chmod +x bin/app-$GOOS-$GOARCH