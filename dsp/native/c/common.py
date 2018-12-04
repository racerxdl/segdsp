#!/usr/bin/env python

import os, subprocess, shutil, re
from templates import getFormattedData, cleanupFileTrash

genasmFolder = "genasm"
outputPrefix = "autogen_"
baseC2GoAsm = "c2goasm -s -c -a {IN} {OUT}"
baseFlags = "-O2 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -S"
baseClang = "clang -D__SUBARCH__={SUBARCH} {BASE_FLAGS} {SUBARCHFLAGS} -masm=intel {CFILE} -o {GENASMFOLDER}/{ASMFILE}"

def initFolders(mainArch):
  if os.path.exists(genasmFolder):
    shutil.rmtree(genasmFolder)

  os.mkdir(genasmFolder)
  os.mkdir(os.path.join("..", mainArch))

def fixFile(filename):
  f = open(filename, "r")
  data = f.read()
  f.close()

  data = cleanupFileTrash(data)
  f = open(filename, "w")
  f.write(data)
  f.close()

def cFilenameToAsm(cFilename, subarch):
  cFilename = os.path.basename(cFilename)
  file, ext = os.path.splitext(cFilename)
  return "%s%s_%s.s" %(outputPrefix, file, subarch)

def cFilenameToArchAsm(cFilename, mainArch, subarch):
  cFilename = os.path.basename(cFilename)
  file, ext = os.path.splitext(cFilename)
  return "%s%s_%s_%s.s" %(outputPrefix, file, subarch, mainArch)

def cFilenameToGoHelper(cFilename, mainArch, subarch):
  cFilename = os.path.basename(cFilename)
  file, ext = os.path.splitext(cFilename)
  return "%s%s_%s_%s.go" %(outputPrefix, file, subarch, mainArch)


def Process(mainArch, source, subarch, outputFolder):
  print("Processing %s => %s" %(source["filename"], subarch["subarchtitle"]))
  goAsmFile = cFilenameToAsm(source["filename"], subarch["name"])
  goAsmFilePath = os.path.join(genasmFolder, goAsmFile)
  goHelperFile = cFilenameToGoHelper(source["filename"], mainArch, subarch["name"])
  goHelperFilePath = os.path.join(genasmFolder, goHelperFile)
  goAsmArchFile = cFilenameToArchAsm(source["filename"], mainArch, subarch["name"])
  goAsmArchFilePath = os.path.join(genasmFolder, goAsmArchFile)

  goFile = getFormattedData(source["cFunction"], {
    "PACKAGE": mainArch,
    "SUBARCH": subarch["name"],
    "SUBARCHTITLE": subarch["subarchtitle"]
  })

  clangLine = baseClang.format(**{
    "SUBARCH": subarch["subarchtitle"],
    "BASE_FLAGS": baseFlags,
    "SUBARCHFLAGS": subarch["flags"],
    "CFILE": source["filename"],
    "ASMFILE": goAsmFile,
    "GENASMFOLDER": genasmFolder,
  })

  print "\nBuilding ==> %s" %goAsmFile
  print "|> %s" %clangLine

  retCode = subprocess.call(clangLine, shell=True)
  print ""

  if retCode != 0:
    print "There was an error compiling %s" %source["filename"]
    exit(retCode)

  # Fix File Crashes for c2goasm
  fixFile(goAsmFilePath)

  # Write Go Helper File
  with open(goHelperFilePath, "w") as f:
    f.write(goFile)

  # Generate Plan9 Assembly
  c2goasmLine = baseC2GoAsm.format(**{
    "IN": goAsmFilePath,
    "OUT": goAsmArchFilePath
  })

  print "\nBuilding ASM ==> %s" %goAsmArchFile
  print "|> %s" %c2goasmLine

  retCode = subprocess.call(c2goasmLine, shell=True)
  print ""

  if retCode != 0:
    print "There was an error compiling %s" %source["filename"]
    exit(retCode)

  print "\nCopying generated files to output at %s" %outputFolder
  shutil.copyfile(goAsmArchFilePath, os.path.join(outputFolder, goAsmArchFile))
  shutil.copyfile(goHelperFilePath, os.path.join(outputFolder, goHelperFile))

def formatFiles():
  subprocess.call("go fmt ../...", shell=True)