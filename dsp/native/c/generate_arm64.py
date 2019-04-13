#!/usr/bin/env python

from common import *

mainArch = "arm64"

outputFolder = "../%s" %mainArch

subarchs = [
  {
    "name": "aarch64",
    "subarchtitle": "arm64",
    "flags": "-I/usr/include/aarch64-linux-gnu/ -target aarch64-none-eabi"
  }
]

sources = [
  {
    "filename": "multiply_conjugate_native.c",
    "function": "MultiplyConjugate",
    "cFunction": "multiplyConjugate",
  },
  {
    "filename": "multiply_conjugate_inline_native.c",
    "function": "MultiplyConjugateInline",
    "cFunction": "multiplyConjugateInline",
  },
  {
    "filename": "dotprod_native_complex.c",
    "function": "DotProductNativeComplex",
    "cFunction": "dotProductComplex",
  },
  {
    "filename": "dotprod_native_float.c",
    "function": "DotProductNativeFloat",
    "cFunction": "dotProductFloat",
  },
  {
    "filename": "dotprod_native_complex_complex.c",
    "function": "DotProductNativeComplexComplex",
    "cFunction": "dotProductComplexComplex",
  },
  {
    "filename": "rotate_complex.c",
    "function": "RotateNativeComplex",
    "cFunction": "rotateComplex",
  },
  {
    "filename": "firfilter.c",
    "function": "FirFilter",
    "cFunction": "firFilter",
  }
]

initFolders()

for source in sources:
  for subarch in subarchs:
    Process(mainArch, source, subarch, outputFolder)

initFolders()
formatFiles()