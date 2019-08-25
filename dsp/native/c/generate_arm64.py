#!/usr/bin/env python

'''
  This is a sketch for the future. Currently c2goasm does not support ARM64 so this is useless.
'''

from common import *

mainArch = "arm64"

outputFolder = "../%s" %mainArch

subarchs = [
  {
    "name": "neon",
    "subarchtitle": "Neon",
    "flags": "-funsafe-math-optimizations -ftree-vectorize -march=armv8-a+simd+fp"
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

print "ARM64 support is still lacking at c2goasm. This won't work"

exit(1)

initFolders(mainArch)

for source in sources:
  for subarch in subarchs:
    Process(mainArch, source, subarch, outputFolder)

initFolders()
formatFiles()