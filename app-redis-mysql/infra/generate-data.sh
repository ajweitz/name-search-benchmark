#!/bin/bash
# to generate names (20k), run: ./generate-list.sh names
# to generate words (500k), run: ./generate-list.sh words

in="../../data/$1.txt"
out='words.txt'
# lowcase='lowcase.txt'
# combined='data.csv'
cp $in $out
# tr '[:upper:]' '[:lower:]' < $out > $lowcase
# paste -d, $out $lowcase > $combined

# rm $out $lowcase