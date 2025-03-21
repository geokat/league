#!/usr/bin/env bash

#set -e

export TARGET=bechallenge;

function runtask() {
    go build -race ./... || return 1
    ./${TARGET}
}

function runwatch() {
    inotifywait -m -e modify -r `pwd` | egrep --line-buffered '\.go$' | while read; do
        # Emacs save can > 1 simultaneous modify events.
        # Avoid restarting > 1 times
        killall -q inotifywait || return
        echo $REPLY
        killall -q ${TARGET}
        echo "Stopping, rebuilding and restarting..."
        runtask &
        sleep 1
    done
}

runtask &
while true; do runwatch; done

