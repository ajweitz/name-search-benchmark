#!/bin/bash
# to generate names (20k), run: ./generate-list.sh names
# to generate words (500k), run: ./generate-list.sh words

in="../../data/$1.txt"
out='words.txt'
temp="temp.txt"
# Ensure each line is at-least 2 chars long
grep -E '^.{2,}$' $in >$out