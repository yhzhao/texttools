#!/usr/bin/env python
"""
Convert between yaml and json
"""

import yaml
import json

import sys

def isJson(l):
    """
    determine input file format based on the first char of l (first line)
    """
    if l is None or len(l) == 0:
      return False
    return l[0]=='{' or l[0] == '['

def yaml2json(filename):
    print("{}".format(next(yaml.load_all(open(filename,"r")))).replace("'","\""))

def json2yaml(filename):
    with open(filename) as jsonfile:
      print(yaml.dump(json.load(jsonfile), allow_unicode=False))

def convert(filename):
    """
    convert file between yaml and json
    """
    with open(filename, 'r') as inputFile:
      firstLine = inputFile.readline()
    if isJson(firstLine):
      json2yaml(filename)
    else:
      yaml2json(filename)

if __name__ == "__main__":
    if sys.version_info[0] < 3:
    	raise Exception("Must use Python 3")
    if len(sys.argv) == 1:
      print("input filename missing")
    else:
       convert(sys.argv[1])

        
