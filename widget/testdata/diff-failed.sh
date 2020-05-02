#!/bin/sh

mkdir diffs
FILES=failed/*
for f in $FILES
do
	n=$(basename -- $f)
	compare -compose src $n $f diffs/$n
done

