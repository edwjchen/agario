import json

CONFIG_PATH = "../peer_to_peer/common/config.json"
BASE_CONFIG_FILE = "base_config.json"

class ConfigGenerator:
    # Takes in list of private ips of EC2 instances + ip of entry server
    def __init__(self, ips, entry_ip):
        self.ips = [ip+":3001" for ip in ips]
        self.entry_ip = entry_ip + ":8080"

    # writes out the config to the default location
    def generate_config(self):
        config = {}
        with open(BASE_CONFIG_FILE) as base:
            config = json.load(base)
        config["ENTRY_SERVER"] = self.entry_ip
        config["REGION_SERVERS"] = self.ips
        with open(CONFIG_PATH, "w") as config_file:
            json.dump(config, config_file, indent=2)
