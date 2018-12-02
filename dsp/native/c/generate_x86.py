#!/usr/bin/env python

from common import *

mainArch = "amd64"

outputFolder = "../%s" %mainArch

subarchs = [
  {
    "name": "avx",
    "subarchtitle": "AVX",
    "flags": "-mavx -mfma"
  },
  {
    "name": "sse2",
    "subarchtitle": "SSE2",
    "flags": "-msse2"
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
  }
]

initFolders()

for source in sources:
  for subarch in subarchs:
    Process(mainArch, source, subarch, outputFolder)

initFolders()
formatFiles()