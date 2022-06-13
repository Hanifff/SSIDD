# SSIDD-server

This project contains the following:
- A Go gRPC server, which consists of:
    - One server (the main.go)
    - proto files
    - Go protobuff generated codes 
    - several rpc methods (implementations in ```./server``` dir.)
- An IPFS cluster, which you can set up with our docker-compose file in ```./ipfs``` dir. 
    - It sets up:
        -   Three Cluster and three Peers
    - A Go IPFS api (implementations in ```./server``` dir.)
- A benchmark module, which consists:
    - A python gRPC stub.
    - Python protobuff generated codes.
    - A python script to benchmark and test gRPC server using ghz module (```bench_ssidd.py```).
    - A script to clean up the network.
    - A shell script to execute tests (```bench.sh```).
- A sample Go gRPC client
- A Dockerfile to set up the server (optional).
- A Hyperledger Fabric SDK

--- 

## Requirements
- cURL - latest version
- Go - 1.17.x
- Docker - 20.10.x
- Docker-compose - 1.29.x
- Pyhton - 3.6.x
- pip3 - 9.0.x
- grpc - 1.45.x
- go-grpc - 1.27.x
- grpcio-tools (python) - v1.46.x
- protobuf - 1.5.x
- Fabric GO SDK - 1.0.x

---

## Set up the server
After successfully installations:<br>
From```./ipfs``` dir run:<br>
```shell
docker-composer up -d
```
From ```./api``` dir run:<br>
```shell
go build
./ssidd
```
Good luck :-)<br>