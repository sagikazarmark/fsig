#!/bin/sh

trap sig 1

sig()
{
    echo "Caught signal"
    exit 0
}

while :
do
    sleep 1
done
