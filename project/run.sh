PROG='/Total Loss:/ {split($2, a, "/"); print a[1]/a[2], a[1], a[2]}'
#'BEGIN {print "Avg.Loss", "TotalLoss", "Count"}
SEP=":"
awk -F$SEP "$PROG" -
