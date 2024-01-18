#!/bin/bash

# Find all .tar.gz files in sharder_blocks and its subdirectories
find sharder-blocks -type f -name "*.tar.gz" -print0 | while IFS= read -r -d '' file; do
    echo "Extracting $file..."
    tar -xzvf "$file" -C /
done