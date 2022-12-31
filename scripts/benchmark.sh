#!/bin/sh


if [ -z $1 ]; then
    echo "No output filename given."
    exit 1
fi

if [ -z $2 ]; then
    n=10
else
    n=$2
fi

logdir="logs"
logfile="$logdir/$1"

mkdir -p "$logdir"
: > "$logfile"

go test -timeout 480m -run='^$' -bench . -count $n > $logfile
