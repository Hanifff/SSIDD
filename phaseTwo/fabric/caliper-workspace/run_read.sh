#!/bin/bash

# This file runs the benchmark module with different setups

# Exit on first error
set -e
starttime=$(date +%s)

clean_up() {
    rm -rf ../test-network/modified_config.json
    rm -rf ../test-network/modified_config.pb
    rm -rf ../test-network/config.pb
    rm -rf ../test-network/config.json
    rm -rf ../test-network/config_update_in_envelope.json
    rm -rf ../test-network/config_update_in_envelope.pb
    rm -rf ../test-network/config_block.json
    rm -rf ../test-network/config_block.pb
}

# reconfigure the network config file with new secret key
python3 -c "import bench_config; bench_config.reconfig_network()"

# benchmark for network 3 org, x clients and x tps
declare -a clients=('5' '10' '20' '50' '100')

for h in "${clients[@]}"; do
    export TPS_Caliper="500"
    export NR_OF_CLINETS="$h"
    echo "Running tests for 3 organisations with tps: $h"
    # rewrite the benchmark configuration
    python3 -c "import bench_config; bench_config.reconfig_benchmark()"
    # sleep for 10 second for configuration to be applied
    sleep 5
    # run the caliper cli
    npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml \
        --caliper-benchconfig benchmarks/bench_read.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled \
        --caliper-report-path results/read/report_DecideRead_FINAL_3org_cli_"$h"__tps_500.html
    # sleep for 5 sec before restarting
    sleep 5
done

# clean up config files
clean_up

# 2 - benchmark by changing block size by varying batch timeout and batch size  (takes more than 4 hours).
declare -a blockSizes=('200' '400' '600' '800' '1000')
declare -a batchTimeOuts=('2s' '4s' '6s' '8s')
export TPS_Caliper='500'
export NR_OF_CLINETS="10"
python3 -c "import bench_config; bench_config.reconfig_network()"
python3 -c "import bench_config; bench_config.reconfig_benchmark()"

for i in "${blockSizes[@]}"; do
    export MAX_BLOCKS="$i"
    for j in "${batchTimeOuts[@]}"; do
        clean_up
        export MY_B_TIMEOUT="$j"
        echo "Running tests for 3 organisations with block size: $h"
        echo "Reconfiguering the channel..."
        # rewrite the channel configuration
        pushd ../accesscontrol
        ./reconf_step1.sh
        popd
        sleep 2
        python3 -c "import bench_config; bench_config.reconfig_channel_blocksize()"
        sleep 2
        echo "Step 2 av reconf.."
        pushd ../accesscontrol
        ./reconf_step2.sh
        popd

        # sleep for 5 second for configuration to be applied
        sleep 5
        echo "Reconfiguerations done."
        echo "Benchmarking..."
        # run the caliper cli
        npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml \
            --caliper-benchconfig benchmarks/bench_read.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled \
            --caliper-report-path results_bs/read/report_read_FINAL_cli_10_bsize_"$i"__btimeout_"$j"_.html
        # sleep for 5 sec before restarting
        sleep 5
    done
done

cat <<EOF

    Total setup execution time : $(($(date +%s) - starttime)) secs ...

    Then, to see the results:
        Open the "./results/report_xx.html" file(s).

EOF
