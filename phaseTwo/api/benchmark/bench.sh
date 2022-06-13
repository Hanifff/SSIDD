#!/bin/bash

set -e
starttime=$(date +%s)

declare -a requests=('2000')
declare -a clients=('5' '10' '20' '50' '100' '500' '1000')

counter=1

# read request
for i in "${requests[@]}"; do
    for j in "${clients[@]}"; do
        python3 bench_ssidd.py ./read_req_"$i"__cli_"$j" "$i" "$j" "read"
    done
    sum=$((counter + $i * 5))
    ./cleanup.go $counter $sum
    counter=$((counter + $i * 5))
done

counter=1
declare -a requests=('500')
declare -a clients=('1' '5' '10')

# write
for i in "${requests[@]}"; do
    for j in "${clients[@]}"; do
        python3 bench_ssidd.py ./bench_results/write/write_data_1Mb_req_"$i"__cli_"$j" "$i" "$j" "write"
    done
    sum=$((counter + $i * 5))
    ./cleanup.go $counter $sum
    counter=$((counter + $i * 5))
done

cat <<EOF
    Total setup execution time : $(($(date +%s) - starttime)) secs ...
    Then, to see the results:
        Open the "./results/report_xx.html" file(s).

EOF
