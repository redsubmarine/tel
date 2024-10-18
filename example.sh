#!/bin/bash

# Set the height of the tree
HEIGHT=5

# Draw the tree
for ((i=1; i<=HEIGHT; i++)); do
    for ((j=i; j<HEIGHT; j++)); do
        echo -n " "
    done

    for ((k=1; k<=2*i-1; k++)); do
        echo -n "*"
    done

    echo ""
    sleep 1
done

for ((i=1; i<=2; i++)); do
    for ((j=1; j<HEIGHT; j++)); do
        echo -n " "
    done
    echo "|"
    sleep 1
done
