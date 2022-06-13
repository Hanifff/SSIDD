# SSIDD's Caliper Benchmark Module

This project contains SSIDD's caliper workspace module. It contains the necessary modules to benchmark the DecideRead and DecideWrite methods of the PDPSC smart contract.

---

## Requirements

Install the Fabric tools and add them to PATH:<br/>

```shell
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.4.0 1.4.8 -s
export PATH=$PATH:$(pwd)/bin
```

Please refer to the main [README.md]("../README.md") file to configuere. You can setup and initialize the network:<br/>

```shell
cd ../accesscontrol
./start.sh
```

In addition, you will need to install:<br>

- Python >= 3.8.x
- pip3

Finally, run the following command to install the necessary python libraries:

```shell
python3 -m pip install -r requirements.txt
```

---

## Benchmark our chaincode

```shell
cd caliper-worksapce
```

Use our script to configuere and run our benchmarking module:<br/>

```shell
./run.sh
```

You can find a `report.html` file in the parent directory.

---
