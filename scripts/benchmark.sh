#!/bin/sh

N=5

logdir="logs"
logfile="$logdir/bench.log"

mkdir -p "$logdir"
: > "$logfile"

for i in $( seq 1 $N ); do
    printf "ITERATION %02d\n" $i
    go test -run=NONE -bench=. -benchmem | tee -a "$logfile"
done

# go test -run=NONE -bench=. > new.txt
