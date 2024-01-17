#!/bin/bash

# Create the base directory
mkdir -p test_data

# Create subdirectories
for i in {0..9}; do
    for j in {a..f}; do
        for k in {0..9}; do
            mkdir -p "test_data/$i/$j/$k/f/1f84345e2ae52f72858d45ef55eb5c8853c85ede165b3edce91d0ba8c065"
        done
    done
done