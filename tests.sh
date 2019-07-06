#!/bin/bash
go build
for f in maps/*
do
	echo "Processing $f file..."
	# take action on each file. $f store current file name
	./N-Puzzle -m $f
done
