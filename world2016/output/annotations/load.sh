#!/bin/bash

FILES=*

REPO="annotate"

pachctl create-repo $REPO
commitid=$(pachctl start-commit $REPO master) 

for f in $FILES
do
	if [ "$f" != "load.sh" ]
	then
		echo $f
		pachctl put-file $REPO $commitid $f -f $f
	fi
done

pachctl finish-commit $REPO $commitid
