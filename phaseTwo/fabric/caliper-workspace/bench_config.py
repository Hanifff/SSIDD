import yaml
import os
import glob
import json


def reconfig_channel_blocksize(file_path: str = '../test-network/modified_config.json') -> None:
    """
        This function reconfiguere Block sizes of the network.

        Args:
            file_path: Path of the copy of config file.
    """
    block_sizes = [200, 400, 600, 800, 1000]
    max_bytes = [524288, 524288*2, 524288*3, 524288*4, 524288*5]
    max_prefered = [103809024, 103809024*2, 103809024*3,103809024*4, 103809024*5]
    with open(file_path, "r+") as config_file:
        configs = json.load(config_file)

    bs = int(os.environ['MAX_BLOCKS'])
    configs["channel_group"]['groups']['Orderer']['values']['BatchSize'][
        'value']['absolute_max_bytes'] = max_prefered[block_sizes.index(bs)]
    configs["channel_group"]['groups']['Orderer']['values']['BatchSize']['value']['max_message_count'] = bs
    configs["channel_group"]['groups']['Orderer']['values']['BatchSize'][
        'value']['preferred_max_bytes'] = max_bytes[block_sizes.index(bs)]

    configs["channel_group"]['groups']['Orderer']['values']['BatchTimeout']['value']['timeout'] = os.environ['MY_B_TIMEOUT']

    with open(file_path, "w") as config_file:
        json.dump(configs, config_file)


def reconfig_network(file_path: str = './networks/networkConfig.yaml') -> None:
    """
        This function modifies the network's CA configurations for the Caliper-workspace.

        Args:
            file_path: Path of the caliper network config file.
    """
    with open(file_path) as f:
        net_config = yaml.safe_load(f)

    sk_path = "../test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/*sk"
    cert_key = "../test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/*.pem"
    cert_keys = []
    sk_file = []
    for file in glob.glob(sk_path):
        sk_file.append(file)
    for file in glob.glob(cert_key):
        cert_keys.append(file)

    net_config['organizations'][0]['identities']['certificates'][0]['clientPrivateKey']['path'] = sk_file[0] if "sk" in sk_file[0] else ""
    net_config['organizations'][0]['identities']['certificates'][0]['clientSignedCert']['path'] = cert_keys[0]

    with open(file_path, 'w') as f:
        yaml.dump(net_config, f)


def reconfig_benchmark(file_path: str = './benchmarks/bench_write.yaml') -> None:
    """
        This function dynamically resets configurations of benchmark module for the caliper-workspace.

        Args:
            file_path: Path to the benchmark configuration file.
    """
    with open(file_path) as f:
        bench_config = yaml.safe_load(f)

    bench_config['test']['workers']['number'] = float(
        os.environ['NR_OF_CLINETS'])
    for round in range(len(bench_config['test']['rounds'])):
        bench_config['test']['rounds'][round]['rateControl']['opts']['tps'] = float(
            os.environ['TPS_Caliper'])

    with open(file_path, 'w') as f:
        yaml.dump(bench_config, f)
