#!/bin/bash

for ((idx=0; idx<=15; idx++))
do
  hex=$(printf "%x" $idx)
  echo "pack $hex"
  mkdir -p sharder-blocks/$hex
  ./packss -path /var/0chain/sharder/hdd/docker.local/sharder1/data/blocks/$hex -dest sharder-blocks/$hex -thread 36

done