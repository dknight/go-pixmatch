#!/bin/sh

# Script generate easier benchmarks.
# Later it is easier to analyze it with benchstat tool.

if [ -z $1 ]; then
    echo "No output filename given."
    exit 1
fi

if [ -z $2 ]; then
    n=10
    echo "Warning! No iteration count given (default $n)."
else
    n=$2
fi

logdir="logs"
logfile="$logdir/$1"

mkdir -p "$logdir"
: > "$logfile"

go test -timeout 5m -run='^$' -bench . -count $n | tee $logfile
