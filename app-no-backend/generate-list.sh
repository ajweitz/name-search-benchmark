#! /bin/bash
in='../names.txt'
out='names.js'
echo names = $(jq -R -s -c 'split("\n")' < $in)>$out