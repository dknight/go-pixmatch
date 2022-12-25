#!/bin/sh

logdir="logs"
logfile="$logdir/bench.log"
N=10

mkdir -p "$logdir"
: > "$logfile"

for i in {1..10}; do
    printf "ITERATION %02d\n" $i
    go test -bench=. | tee -a "$logfile"
done
