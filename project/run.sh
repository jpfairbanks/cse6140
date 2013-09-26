#!/usr/bin/env bash
PROG='/Total Loss:/ {split($2, a, "/"); print a[1]/a[2], a[1], a[2]}'
#'BEGIN {print "Avg.Loss", "TotalLoss", "Count"}
SEP=":"
for (( i=1; i<=100; i+=10 )); do
    ./project -depth $i -width 700 -efactor 8 |\
        awk -F$SEP "$PROG" -
done
