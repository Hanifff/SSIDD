# Phase one components

This project contains phase components of SSIDD:
The ```agents``` directory contains three ASP .Net Core 3.1 web applications:
  - A TA (trust anchor), which contains:
      - An UI
      - An API 
      - Aries agent 
  - A client (servs clients), which contains:
      - An UI
      - An API
      - Aries agent 
  - A gatekeeper (policy enforcement point of access controller), which contains:
      - An UI
      - An API
      - Aries agent 
      - a gRPC proto and client module

<br>
A docker-compose file to set up the network. It sets up:
  - Four docker containers for an Hyperledger Indy pool
  - One docker container for TA 
  - One docker container for Client
  - One docker container for gatekeeper (pep)
<br>
A docker-compose file with mounted directory as volume to implement and see the changes in real time. It sets up:
  - Four docker containers for an Hyperledger Indy pool
  - One docker container for TA 
  - One docker container for Client
  - One docker container for gatekeeper (pep)

Docker file can be found in ```./docker``` directory.

All materials in this project are implemented using [Hyperledger Aries .Net framework](https://github.com/hyperledger/aries-framework-dotnet) as baseline.

---

## Requirements
In order to setup the project the following tools should be installed:

- cURL - latest version
- .Net core - 3.1
- Node js - 16.14.x
- npm - 8.3.x
- Docker - 20.10.x
- Docker-compose - 1.29.x
- AgentFramework - 4.0.x
- Hyperledger.Aries.AspNetCore - 1.6.x
- libindy - latest version
- Hyperledger Indy CLI - 1.16.x

---

## Running the projects
After successfully installations, run one of the two docker-compose files.<br>
Here we go, you have everything up and running :-)<br>