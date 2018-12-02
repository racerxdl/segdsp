#!/bin/bash

function genFile {
  INST="$1"
  ARCH="$2"
  PACK="$3"
cat << EOF >> dotprod_native_${ARCH}_${INST}.go
//+build !noasm
//+build !appengine

package ${PACK}

import (
  "unsafe"
)

//go:noescape
func _dotProductComplex_${INST}(result, input, taps unsafe.Pointer, length uint)

//go:noescape
func _dotProductFloat_${INST}(result, input, taps unsafe.Pointer, length uint)

func DotProductNativeFloat_${INST}(input []float32, taps []float32) float32 {
  var res = make([]float32, 1)

  var resPtr = unsafe.Pointer(&res[0])
  var inputPtr = unsafe.Pointer(&input[0])
  var tapsPtr = unsafe.Pointer(&taps[0])
  var cLen = uint(len(taps))

  _dotProductFloat_${INST}(resPtr, inputPtr, tapsPtr, cLen)

  return res[0]
}
func DotProductNativeComplex_${INST}(input []complex64, taps []float32) complex64 {
  var res = make([]float32, 2)

  var resPtr = unsafe.Pointer(&res[0])
  var inputPtr = unsafe.Pointer(&input[0])
  var tapsPtr = unsafe.Pointer(&taps[0])
  var cLen = uint(len(taps))

  _dotProductComplex_${INST}(resPtr, inputPtr, tapsPtr, cLen)

  return complex(res[0], res[1])
}
EOF
}

mkdir -p genasm
rm -f genasm/*

BASE_FLAGS="-fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -S"

# Generate x64 SIMD
clang -D__SUBARCH__=arm64 -O2 -target armv8 -mfpu=neon -mfloat-abi=hard ${BASE_FLAGS} dotprod_native.c -o genasm/dotprod_native_neon.s

mkdir -p ../arm64
cd genasm

for i in `ls dotprod_native_*.s`
do
  echo "----- Processing $i -----"
  instruction=`echo $i | sed 's/dotprod_native_\(.*\).s/\1/g'`
  goasmfile=`echo $i | sed 's/dotprod_native_\(.*\).s/dotprod_native_arm64_\1.s/g'`
  genFile $instruction arm64 arm64
  c2goasm -a $i $goasmfile
  mv $goasmfile ../../arm64
  echo ""
done

mv *.go ../../arm64
