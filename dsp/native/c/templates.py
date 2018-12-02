#!/usr/bin/env python
import re

def loadTPL(funcName):
  f = open("tpls/%s.go_tpl" % funcName)
  tpl = f.read()
  f.close()
  return tpl

def getFormattedData(funcName, args):
  tplData = loadTPL(funcName)
  return tplData.format(**args)

trashR =[
  re.compile(r".*Lfunc_end0.*"),
  re.compile(r".*.size.*"),
  re.compile(r".*.ident.*"),
  re.compile(r".*section.*"),
]

def cleanupFileTrash(data):
  o = ""
  for line in data.split("\n"):
    for reg in trashR:
      line = re.sub(reg, "", line)
    o += line + "\n"
  return o.strip()