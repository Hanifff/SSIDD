"""
This script models the gatekeeper communication with grpc server integrated with IPFS and HLF.
Script can be executed as: python3 bench_ssidd.py file_name  nr_of_requests(def: 200) nr_of_clients type(read/write)
"""
import json
import string
import random
from math import fabs
import sys
import os
import subprocess
from time import sleep, time
import ssidd_pb2
import datetime
from google.protobuf.timestamp_pb2 import Timestamp
from google.protobuf.struct_pb2 import Value


CURRENT_DIR = os.path.dirname(os.path.abspath(__file__))
R_PROTO = "../protos/ssidd.proto"
GHZ = "ghz"


def dataGenerator(N: int) -> str:
    """ generates n bytes data. """
    return ''.join(random.choice(string.ascii_lowercase) for _ in range(N)).lower()

             
def bench_read(nr_of_benchs: int, nr_of_clis: int, nr_of_msgs: int, bench_name: str, port: str):
    """ Benchmarks the Read rpc method. 
        args: 
            nr_of_bench: number of benchmark rounds.
            nr_of_cli: number of conccurent clients.
            nr_of_messages: number of message to be sent in each round.
            port: gRPC server's address:port number.
    """
    t = datetime.datetime.now().timestamp()
    seconds = int(t)
    nanos = int(t % 1 * 1e9)
    attr_cli = {"status": Value(bool_value=True), "expiration": Value(string_value="02-Jan-2026"),
                "organisationid": Value(string_value="spain01")}
    read_message= ssidd_pb2.ReadRequest(ClientDID= "did:client:123456789abcdefghigklmn",ResourceID= "r004", 
                                ClientAttributes=attr_cli,
                                Timestamp= Timestamp(seconds=seconds, nanos=nanos),
                                IsPolicy=False)
    bin_fpath = "./bench_results/read_bench.bin"
    serialized_msg = read_message.SerializeToString()
    bin_file = open(bin_fpath, "wb")
    bin_file.write(serialized_msg)
    bin_file.close()
    
    for i in range(nr_of_benchs):
        execute([GHZ, "--insecure", "--proto", R_PROTO, "--call",
            "ssidd.Ssidd.Read",  "--cpus", str(8), "-c", str(nr_of_clis), "--connections",str(1), "-n", str(nr_of_msgs),
            "-B", bin_fpath,"-t", "400s", "-o",
            "(" + str(i+1) + ")" + bench_name + str(nr_of_msgs) + ".json", "-O", "pretty", port])



def bench_write(nr_of_benchs: int, nr_of_clis: int, nr_of_msgs: int, bench_name: str, port: str, data:str):
    """ Benchmarks the Write rpc method. 
        args: 
            nr_of_bench: number of benchmark rounds.
            nr_of_cli: number of conccurent clients.
            nr_of_messages: number of message to be sent in each round.
            port: gRPC server's address:port number.
            data: data to be sent along with WriteRequest message.
    """
    attr_cli = {"status": Value(bool_value=True), "expiration": Value(string_value="02-Jan-2026"),
                "organisationid": Value(string_value="spain01")}
    attr_res = {"any-key":Value(string_value="any-value")}
    for i in range(nr_of_benchs): 
        write_message = ssidd_pb2.WriteRequest(ClientDID= "did:client:123456789abcdefghigklmn", 
                                            ResourceID=dataGenerator(4), PolicyID ="p001",OwnerOrgID="spain01",
                                            ClientAttributes=attr_cli, Data= bytes(data, encoding='utf-8'),
                                            OptionalResAttributes= attr_res)
        
        bin_fpath = "./bench_results/write_bench.bin"
        serialized_msg = write_message.SerializeToString()
        bin_file = open(bin_fpath, "wb")
        bin_file.write(serialized_msg)
        bin_file.close()
        execute([GHZ, "--insecure", "--proto", R_PROTO, "--call",
                    "ssidd.Ssidd.Write", "-c", str(
                        nr_of_clis), "--connections",str(nr_of_clis), "-n", str(nr_of_msgs),
                    "-B", bin_fpath,"-t", "400s",  "-o",
                    "(" + str(i+1) + ")" + bench_name + str(nr_of_msgs) + ".json", "-O", "pretty", port])
    
    

def extract_results(file_path: str):
    """ Extracts benchmark results from json files.
        args:
            file_path: a new file with average of results.
    """
    total_latency = 0
    total_throughput = 0
    total_resp_time = 0
    result = {}
    result["benchmarks"] = []
    result["final_Avg."] = []

    obj = os.scandir(CURRENT_DIR)
    for entry in obj:
        if entry.is_file():
            if entry.name.endswith(".json"):
                myfile = os.path.join(CURRENT_DIR, entry.name)
                with open(myfile, "r") as f:
                    data = json.load(f)
                    details = data["details"]
                    sum_latency = 0
                    for listing in details:
                        sum_latency += listing["latency"]
                    avg = sum_latency / len(details)

                    result["benchmarks"].append({
                        "file_name": entry.name.replace(".json", ""),
                        "avg_latency(ms)": float(avg * (10 ** -6)), 
                        "avg_respons_time(ms)": data["average"] * (10 ** -6),
                        "requests/sec(throughput)": data["rps"],
                    })
                    total_latency += float(avg * (10 ** -6))
                    total_resp_time += data["average"] * (10 ** -6)
                    total_throughput += data["rps"]

    number_bs = len(result["benchmarks"])
    result["final_Avg."].append({
        "latency(ms)": total_latency / number_bs,
        "respons_time(ms)": total_resp_time / number_bs,
        "requests/sec(throughput)": total_throughput / number_bs,
    })

    newfile = os.path.join(CURRENT_DIR, file_path + ".json")
    with open(newfile, "w") as outfile:
        json.dump(result, outfile, indent=3, sort_keys=True)


def cleanUp(name: str):
    obj = os.scandir(CURRENT_DIR)
    for entry in obj:
        if name in entry.name:
            os.remove(entry.name)


def execute(cmd):
    p = subprocess.Popen(cmd)
    p.wait()

def main():
    result_file = sys.argv[1]
    nr_of_msgs = sys.argv[2]
    nr_of_clis = sys.argv[3]
    req_type = sys.argv[4]
    addrDocker = "10.0.0.14:8081"
    addrLocal = "127.0.0.1:8081"
    
    if req_type == "read":
        bench_read(10, nr_of_clis, nr_of_msgs, os.path.basename(
        CURRENT_DIR),addrLocal)
    if req_type == "write":
        data = dataGenerator(500000000) # x kb/mb/gb data
        bench_write(1, nr_of_clis, nr_of_msgs,os.path.basename(
        CURRENT_DIR), addrLocal, data)
   
    extract_results(result_file)
    print("os.path.basename(CURRENT_DIR) ", os.path.basename(CURRENT_DIR))
    cleanUp(os.path.basename(CURRENT_DIR))


if __name__ == "__main__":
    main()