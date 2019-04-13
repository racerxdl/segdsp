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
    "filename": "add_complex_complex_vectors.c",
    "function": "AddComplexComplexVectors",
    "cFunction": "addComplexComplexVectors",
  },
  {
    "filename": "subtract_complex_complex_vectors.c",
    "function": "SubtractComplexComplexVectors",
    "cFunction": "subtractComplexComplexVectors",
  },
  {
    "filename": "multiply_complex_complex_vectors.c",
    "function": "MultiplyComplexComplexVectors",
    "cFunction": "multiplyComplexComplexVectors",
  },
  {
    "filename": "divide_complex_complex_vectors.c",
    "function": "DivideComplexComplexVectors",
    "cFunction": "divideComplexComplexVectors",
  },
  {
    "filename": "add_float_float_vectors.c",
    "function": "AddFloatFloatVectors",
    "cFunction": "addFloatFloatVectors",
  },
  {
    "filename": "subtract_float_float_vectors.c",
    "function": "SubtractFloatFloatVectors",
    "cFunction": "subtractFloatFloatVectors",
  },
  {
    "filename": "multiply_float_float_vectors.c",
    "function": "MultiplyFloatFloatVectors",
    "cFunction": "multiplyFloatFloatVectors",
  },
  {
    "filename": "divide_float_float_vectors.c",
    "function": "DivideFloatFloatVectors",
    "cFunction": "divideFloatFloatVectors",
  },
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