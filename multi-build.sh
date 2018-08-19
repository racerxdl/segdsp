#!/bin/bash

REPO="racerxdl/segdsp"

echo "Building ${REPO}:latest"
docker build -t racerxdl/segdsp .

echo "Building arm32v6 build locally (bug in go compiler inside alpine)"

go get -v

CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o segdsp_worker_arm32v6

echo "Downloading qemu-user for archs"

for target_arch in aarch64 arm x86_64
do
  echo "Downloading ${target_arch} qemu userspace"
  wget -N https://github.com/multiarch/qemu-user-static/releases/download/v2.12.0-1/x86_64_qemu-${target_arch}-static.tar.gz
  tar -xvf x86_64_qemu-${target_arch}-static.tar.gz
done

for arch in amd64 arm32v6 arm64v8; do
  echo "Building ${REPO}:${arch}-latest"
  docker build -f Dockerfile.${arch} -t ${REPO}:${arch}-latest .
done


echo "Pushing ${REPO}:latest"
docker push racerxdl/segdsp

for arch in amd64 arm32v6 arm64v8; do
  echo "Pushing ${REPO}:${arch}-latest"
  docker push ${REPO}:${arch}-latest
done

docker images racerxdl/segdsp