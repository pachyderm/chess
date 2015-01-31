# Chess repo
This repo contains code to perform distributed analysis on chess games (.pgn files).
It was used to create this [article](https://medium.com/@jdoliner/when-grandmasters-blunder-a819860b883d).

# Overview
## data
The `data` directory contains the raw data that was used to create this demo.
It also includes job specifications (`chessMap` and `chessReduce`)  which will
launch the pfs jobs that actually perform the analysis and a script which will
send all of these to a running pfs instance using curl and kickoff the pipeline
with a commit.

## output
The `output` contains the results of running the job that we got from pfs
(`output/blunders`) as well as a version of that data bucketed by rating
`output/bucketed`

## `map`/`reduce`
`map` and `reduce` contain the actual implementation of the jobs used to do the
processing.

## Dockerfile
`Dockerfile` can be used to build the repo in to a container which you can give
to pfs. This container already exists in docker hub at `pachyderm/chess` but
you can use the Dockerfile to create a container with your own jobs. To use
your new jobs you'll need to change the image used in data/chessMap and
data/chessReduce.
