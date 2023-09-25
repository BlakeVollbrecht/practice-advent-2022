#! /bin/bash

{
  # printf total calories for each elf
  sum=0

  while IFS="" read -r line; do
    if [ -z "$line" ]; then
      printf "%s\n" "$sum"
      sum=0
    else
      sum=$((sum + line))
    fi
  done < ./inputs/1.txt
} |
sort -nr |
head --lines=3 |
awk '{n += $1}; END{print n}' # sum top 3 elves' calories
