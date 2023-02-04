#!/bin/bash
# to generate names (20k), run: ./generate-list.sh names
# to generate words (500k), run: ./generate-list.sh words

in="../../data/$1.txt"
out='words.txt'
lowcase='lowcase.txt'
cp $in $out
tr '[:upper:]' '[:lower:]' < $out > $lowcase