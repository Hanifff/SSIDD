# Self-Sovereign Identity based, Dynamic, and Decentralized Access Control (SSIDD)

This repository contains implementations of the SSIDD, a novel access control system that ensures user privacy, access policy flexibility, and trust.</br>

The work was done as a master thesis project at the University of Stavanger.<br>

## Contains
The repository contains several projects that work in parallel to provide access control services.<br>
Those are:
- Phase one components:
    - Hyperledger Indy ledger
    - Three ASP .Net applications that implement Hyperledger Aries and gRPC.
        - A TA (trust anchor)
        - A client (serves clients)
        - A gatekeeper (policy enforcement point of access controller)

- Phase two components:
    - SSIDD-server:
        - Go gRPC
        - Python gRPC
        - benchmark module
        - IPFS
        - Fabric SDK
    - Hyperledger Fabric:
        - Smart contracts (chaincode) and unit tests
        - caliper workspace module

Please refer to each folder for more details.<br>

## Set up the access controller
In order to set up the entire application, you need to follow the description details of each component in their respective README files.<br>
Additionally, you may also need to visit the official documentation for:
- Hyperledger Fabric, Hyperledger Indy, or Hyperledger Aries
- The GO, Python, C#, Node js programming languages
- Docker, 
- The Go gRPC, .Net gRPC, or Python gRPC
- IPFS and Go IPFS API

### Best of luck!