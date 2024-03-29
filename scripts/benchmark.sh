#!/bin/sh

# Script generate easier benchmarks.
# Later it is easier to analyze it with benchstat tool.

logdir="logs"
logfile="$logdir/$1"

if [ -z $1 ]; then
    echo "No output file name given."
    exit 1
fi

if [ -z $2 ]; then
    n=10
    echo "Warning! No iteration count given (default $n)."
else
    n=$2
fi

mkdir -p "$logdir"
: > "$logfile"

go test -timeout 5m -run='^$' -bench . -count $n | tee $logfile
