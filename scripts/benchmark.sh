#!/bin/sh

n=$2

if [ -z $N ]; then
    n=20
fi

if [ -z $1 ]; then
    echo "No output filename given."
    exit 1
fi

logdir="logs"
logfile="$logdir/$1"

mkdir -p "$logdir"
: > "$logfile"

go test  -timeout 480m -run='^$' -bench=. -count=$n > $logfile
