# SSIDD Chaincode 
 
This project contains the SSIDD's smart contracts and a Caliper workspace module. 

## Requirements

---

This project is implemented in [Go](https://golang.org/)<br/>
You need to have a Hyperledger Fabric network running with at least three peer organizations and an ordering service.<br/>
However, fabric provides a [test network](https://hyperledger-fabric.readthedocs.io/en/release-2.2/test_network.html) which you can run locally and has all required configurations. In the following, we list the requirements for compiling this application:<br>

- Fabric — version 2.4
- cURL — latest version
- Go — version 1.17.x
- Docker Compose — version 1.29.x
- Hyperledger Fabric test network<br/>


---
##  Start the test network, configure a channel, deploy the chaincode, and start the application server


Alternatively, you can use the script shell file to setup all configurations mentioned above:<br>

```shell
cd ../accesscontrol
./start.sh
```

You can shot down the network and delete all its dependencies using the following command:<br>

```shell
./shotdown.sh
```
---
#### Good luck :-)