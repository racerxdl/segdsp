#!/bin/bash

ARCH=`uname -m`

if [ "${ARCH}" = "x86_64" ] || [ "${ARCH}" = "i386" ]
then
  echo "------------- x86 ------------"
  ./generate_amd64.sh
  # ./generate_x86.sh
elif [ "${ARCH}" = "aarch64" ] || [ "${ARCH}" = "arm64" ]
then
  echo "------------- arm64 -----------"
  ./generate_arm64.sh
fi
