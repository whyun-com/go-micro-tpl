#!/usr/bin/python
# -*- coding: UTF-8 -*-

import os, subprocess, sys
tagName = os.getenv('CI_COMMIT_TAG')
ret = os.system('git fetch && git checkout ' + tagName)
if ret != 0:
    sys.exit(ret)

checkBranch = os.getenv('CHECK_BRANCH')
if checkBranch == None or checkBranch == '__NONE':
    print('not need to check branch')
    sys.exit(0)

output = subprocess.check_output('git branch -a --contains tags/' + tagName, shell=True)

lines = output.split('\n')

isValidTag = False
for line in lines:
    name = line.replace('*','').replace(' ','')
    if name == checkBranch or name == 'remotes/origin/' + checkBranch:
        isValidTag = True
        break

if isValidTag:
    sys.exit(0)
else:
    print('invalid tag '+ tagName +' for ' + checkBranch)
    sys.exit(255)