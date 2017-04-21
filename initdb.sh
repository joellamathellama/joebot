#!/bin/sh

for f in [0-9]*.sql
do
    echo "Attempting to load $f."
    psql "$@" -f "$f"
    if [ $? -ne 0 ]
    then
        echo "Failed to load $f."
        exit
    else
        echo "Loaded $f."
    fi
    echo
done
