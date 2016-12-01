#!/bin/bash

FILES=*.pgn

pachctl create-repo chess
commitid=$(pachctl start-commit chess master) 

for f in $FILES
do
	pachctl put-file chess $commitid $f -f $f
done

pachctl finish-commit chess $commitid
