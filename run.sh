for ((idx=0; idx<=15; idx++))
do
  hex=$(printf "%x" $idx)
  echo "pack $hex"
  packss -path /var/0chain/sharder/hdd/docker.local/sharder1/data/block/$hex -dest $hex

done