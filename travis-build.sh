#!/bin/bash

TAG=`git describe --exact-match --tags HEAD`

if [ $? -eq 0 ];
then
  echo "Releasing for tag ${TAG}"
  echo "Downloading deps"
  go get -v
  echo "Dowloading gox for multi-arch"
  go get github.com/mitchellh/gox
  mkdir out
  mkdir bins
  echo "Multi-arch build"
  gox -output "out/{{.OS}}-{{.Arch}}/{{.Dir}}" -arch="arm arm64 386 amd64" -os="windows linux"
  cd out
  for i in *
  do
    echo "Zipping segdsp-${i}.zip"
    cp -Rv ../content $i/
    zip -r ../bins/segdsp-$i.zip $i/*
  done
  cd ..
  ls -la bins
else
  echo "No tags for current commit. Skipping releases."
fi
