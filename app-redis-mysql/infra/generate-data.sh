#!/bin/bash
# to generate names (20k), run: ./generate-list.sh names
# to generate words (500k), run: ./generate-list.sh words

in="../../data/$1.txt"
out='words.txt'
cp $in $out