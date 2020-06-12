# run wait_start then
# run setup first, then run refresh to update all nodes with the new config file

import boto3
import os
import time
import argparse
import pprint
import multiprocessing
import paramiko
import subprocess
import sys
from config_generator import ConfigGenerator

parser = argparse.ArgumentParser(description="This program autocreates, starts, and stop Amazon EC2 Instancess.")
parser.add_argument('-create', default=0, help="Number of new instances to create")
parser.add_argument('-start', default=0, help="Number of new instances to start")
parser.add_argument('-stop', default=0, help="Number of new instances to stop")
parser.add_argument('-terminate', default=0, help="Number of new instances to terminate")
parser.add_argument('-stats', default=False, help="Number of new instances to terminate")
parser.add_argument('-run', default=True, help="Command to run on running instances")
parser.add_argument('-setup', default=False, help="Setup instances")
parser.add_argument('-hostname', default=False, help="Print all hostnames of running instances")
parser.add_argument('-verify', default=False, help="verify setup")

ENTRY_NAME = ""
SERVER_NAMES = []
EXPERIMENT_NAME= ""

ec2_resource = boto3.resource('ec2',
    aws_access_key_id=os.environ['AWS_ACCESS_KEY_ID'],
    aws_secret_access_key=os.environ['AWS_SECRET_ACCESS_KEY'],
    region_name='us-west-1')

def create_instances(num):
    instance = ec2_resource.create_instances(ImageId='ami-0318e6f2445586bd7',
        InstanceType='t2.micro',
        KeyName='the-key-to-her-heart',
        MinCount=1,
        MaxCount=num,
        Monitoring={'Enabled': False},
        SecurityGroups=['agario']
    )

def start_instances(num):
    stopped_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['stopped']}]))
    count = 0
    num = min(num, len(stopped_instances))
    for i in range(num):
        instance = stopped_instances[i]
        ec2_resource.instances.filter(InstanceIds=[instance.id]).start()
        count += 1
    print("Started {} instances".format(count))


def stop_instances(num):
    running_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]))
    count = 0
    num = min(num, len(running_instances))
    for i in range(num):
        instance = running_instances[i]
        print('Stopping ', instance.public_dns_name)
        ec2_resource.instances.filter(InstanceIds=[instance.id]).stop()
        count += 1
    print("Stopped {} instances".format(count))


def terminate_instances(num):
    instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running', 'stopped']}]))
    count = 0
    num = min(num, len(instances))
    for i in range(num):
        instance = instances[i]
        count += 1
        ec2_resource.instances.filter(InstanceIds=[instance.id]).terminate()
    print("Terminated {} instances".format(count))

def get_stats():
    stats = {}
    stats['total'] = len(list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running', 'stopped', 'pending', 'stopping']}])))
    stats['running'] = len(list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}])))
    stats['stopped'] = len(list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['stopped']}])))
    stats['pending'] = len(list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['pending']}])))
    stats['stopping'] = len(list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['stopping']}])))
    pp = pprint.PrettyPrinter(indent=4)
    pp.pprint(stats)

# waits for n servers to start
def wait_start(num_servers):
    while True:
        running_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]))
        print(len(running_instances), "servers are up.")
        if len(running_instances) >= num_servers:
            break
        time.sleep(1)

# generates config and setups which servers to run
def setup_experiment(server_num):
    global ENTRY_NAME
    global SERVER_NAMES
    running_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]))
    if (server_num > len(running_instances)):
        print("RUN STUFF PROPERLY")
        return

    entry_instance = running_instances.pop()
    ENTRY_NAME = entry_instance.public_dns_name
    entry_private_ip = entry_instance.private_ip_address

    server_instances = running_instances
    SERVER_NAMES = [instance.public_dns_name for instance in server_instances]
    server_private_ips = [instance.private_ip_address for instance in server_instances]

    cg = ConfigGenerator(server_private_ips, entry_private_ip)
    cg.generate_config()

    commit_msg = "Config update for " + EXPERIMENT_NAME
    subprocess.call('git add ../peer_to_peer/common/*.json; git commit -m "'+commit_msg+'"; git push', shell=True)
    time.sleep(1)

    refresh_instances()

def get_hostnames():
    pp = pprint.PrettyPrinter(indent=4)
    pp.pprint([instance.public_dns_name for instance in ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}])])

def worker(ip):
    key = paramiko.RSAKey.from_private_key_file("the-key-to-her-heart.pem")
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    client.connect(hostname=ip, username="ubuntu", pkey=key)
    stdin, stdout, stderr = client.exec_command('rm -rf agario > /dev/null; git clone https://github.com/edwjchen/agario.git; export GOPATH=/home/ubuntu/agario; export PATH=$PATH:/usr/local/go/bin:/home/ubuntu/.local/bin; source .bashrc; bash agario/src/setup/setup.sh')
    stdin.flush()

    if stdout.channel.recv_exit_status():
        print(ip, " failed clone")

    # stdin, stdout, stderr = client.exec_command('')
    # stdin.flush()
    # if stdout.channel.recv_exit_status():
    #     print(ip, " failed setup")

    client.close()

def verify_instance_setup():
    running_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]))
    running_instance_name = [instance.public_dns_name for instance in running_instances]
    def verify(idx, name):
        print(idx, " check: ", name)
        key = paramiko.RSAKey.from_private_key_file("the-key-to-her-heart.pem")
        client = paramiko.SSHClient()
        client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

        client.connect(hostname=name, username="ubuntu", pkey=key)
        stdin, stdout, stderr = client.exec_command('cd agario/src/github.com; ls')
        stdin.flush()

        if stdout.channel.recv_exit_status():
            print(name, " failed setup")
    for idx, i in enumerate(running_instance_name):
        verify(idx, i)

def setup_instances(num):
    running_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]))
    if len(running_instances) < num:
        print("Not all instances are up yet!")
        return

    running_instance_ips = [instance.public_dns_name for instance in running_instances]
    pool = multiprocessing.Pool(len(running_instance_ips))
    pool.map(worker, running_instance_ips)

def refresh(name):
    print("Refreshing ", name)
    key = paramiko.RSAKey.from_private_key_file("the-key-to-her-heart.pem")
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    client.connect(hostname=name, username="ubuntu", pkey=key)
    stdin, stdout, stderr = client.exec_command('cd agario; git pull')
    stdin.flush()

    if stdout.channel.recv_exit_status():
        print(name, " failed clone")
    client.close()


def refresh_instances():
    print("Commencing instance refresh")
    running_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]))
    running_instance_ips = [instance.public_dns_name for instance in running_instances]
    pool = multiprocessing.Pool(len(running_instance_ips))
    pool.map(refresh, running_instance_ips)

def start_single_server(dns_name):
    print("Starting single server on", dns_name)

    key = paramiko.RSAKey.from_private_key_file("the-key-to-her-heart.pem")
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    client.connect(hostname=dns_name, username="ubuntu", pkey=key)
    stdin, stdout, stderr = client.exec_command('export GOPATH=/home/ubuntu/agario; /usr/local/go/bin/go run agario/src/peer_to_peer/main.go')
    stdin.flush()

    client.close()

def start_servers(num):
    global ENTRY_NAME
    running_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]))
    if len(running_instances) - 1 < num:
        print("!!Error!!: Need at least", num+1, "instances running!")
        return

    if ENTRY_NAME == "":
        print("!!Error!!: Start entry first!")
        return

    running_instances_names = [instance.public_dns_name for instance in running_instances]
    try:
        running_instances_names.remove(ENTRY_NAME)
    except:
        print("!!Error!!: COME ON MAN")
    
    # TODO add filtering here to not start on servers already running servers
    pool = multiprocessing.Pool(num)
    pool.map(start_single_server, running_instances_names[:num])

def start_entry():
    global ENTRY_NAME

    key = paramiko.RSAKey.from_private_key_file("the-key-to-her-heart.pem")
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    print("Starting entryserver on", ENTRY_NAME)
    client.connect(hostname=ENTRY_NAME, username="ubuntu", pkey=key)
    stdin, stdout, stderr = client.exec_command('export GOPATH=/home/ubuntu/agario; /usr/local/go/bin/go run agario/src/peer_to_peer/entryserver.go')
    stdin.flush()

    time.sleep(5)

    _, stdout, _ = client.exec_command('lsof -i :8080')
    if stdout.channel.recv_exit_status():
        print(ENTRY_NAME, "failed to start entryserver")
    else:
        print(ENTRY_NAME, "successfully started entryserver")

    client.close()

# type is name of process either ['client', 'entry', 'server']
def killall(dns_name, ptype):
    if ptype == 'entry':
        cmd = "killall -2 entryserver"
    elif ptype == 'server':
        cmd = "killall -2 main"
    else:
        cmd = "killall -2 Python"
    
    print("Killing ", ptype, " on ", dns_name)

    key = paramiko.RSAKey.from_private_key_file("the-key-to-her-heart.pem")
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    client.connect(hostname=dns_name, username="ubuntu", pkey=key)
    stdin, stdout, stderr = client.exec_command(cmd)
    stdin.flush()
    time.sleep(2)
 
    if stdout.channel.recv_exit_status():
        print(dns_name, "failed to kill")
    else:
        print("Killed", ptype, "successfully")
    client.close()

def get_logs(dns_name):
    key = paramiko.RSAKey.from_private_key_file("the-key-to-her-heart.pem")
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    client.connect(hostname=dns_name, username="ubuntu", pkey=key)
    stdin, stdout, stderr = client.exec_command('scp agario/src/peer_to_peer/client/logs/* ../data/experiment_{};'.format(EXPERIMENT_NAME))
    stdin.flush()
    time.sleep(5)
    if stdout.channel.recv_exit_status():
        print(dns_name, " failed to scp")
    client.close()

# instance_ids = []

# ec2_client = boto3.client('ec2')
# response = ec2_client.describe_instances()
# for reservation in response["Reservations"]:
#     for instance in reservation["Instances"]:
#             print(instance["InstanceId"])
#             print(instance)
#             print()
#             instance_ids.append(instance["InstanceId"])


# stopped_instances = ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['stopped', 'stopping']}])
# for instance in stopped_instances:
#     print(instance.id, instance.instance_type)
#     print(instance.public_ip_address)
#     print()

if __name__ == '__main__':
    if (len(sys.argv) == 1):
        print("Missing experiment name!")
        sys.exit(0)

    EXPERIMENT_NAME = sys.argv[1]
    # YES I KNOW makedirs EXISTS BUT WE GOT NO TIME TO FIGURE OUT HOW TO USE IT THANKS - Godwin '20
    if not os.path.exists('../data'):
        os.mkdir('../data')
    if not os.path.exists('../data/'+EXPERIMENT_NAME):
        os.mkdir('../data/'+EXPERIMENT_NAME)

    print("Welcome to the agar.io experimental CLI.")
    print("You are running experiment " + EXPERIMENT_NAME)
    last_cmd = ""
    while True:
        cmds = input('> ').split(" ")
        cmd_type = cmds[0]

        # press enter to redo
        if cmd_type == "" and last_cmd != "":
            cmd_type = last_cmd
        else:
            last_cmd = cmd_type

        if cmd_type == "help":
            print("Not again... oh well here you go\n")
            print("stats hostnames verify create start stop")
            print("terminate setup wait start_entry kill_entry")
            print("start_servers start_clients kill_server_client")
            print("data refresh wait_start")
        elif cmd_type == "stats":
            get_stats()
        elif cmd_type == "hostnames":
            get_hostnames()
        elif cmd_type == "verify":
            verify_instance_setup()
        elif cmd_type == "create":
            create_instances(int(cmds[1]))
        elif cmd_type == "start":
            start_instances(int(cmds[1]))
        elif cmd_type == "stop":
            stop_instances(int(cmds[1]))
        elif cmd_type == "terminate":
            terminate_instances(int(cmds[1]))
        elif cmd_type == "setup":
            setup_instances(int(cmds[1]))
        elif cmd_type == "setup_experiment":
            setup_experiment(int(cmds[1]))
        elif cmd_type == "wait":
            time.sleep(int(cmds[1]))
        elif cmd_type == "refresh":
            refresh_instances()
        elif cmd_type == "wait_start":
            wait_start(int(cmds[1]))
        elif cmd_type == "start_entry":
            start_entry()
        elif cmd_type == "kill_entry":
            killall(ENTRY_NAME, 'entry')
        elif cmd_type == "start_servers": # starts some number of player servers
            num_to_start = int(cmds[1])
            pass
        elif cmd_type == "start_clients": # starts clients on some number of servers that are already running player servers
            pass
        elif cmd_type == "kill_server_client": # kills some number of player servers along with their clients
            num_to_kill = int(cmds[1])
            pass
        elif cmd_type == "data": # run scp on all servers that are running
            pass
        elif cmd_type in ['quit', 'q']:
            sys.exit(0)
        else:
            print("GIT GUD")
        # elif cmd_type == "start_entry":

# if __name__ == '__main__':
#     args = parser.parse_args()
#     if bool(args.stats):
#         get_stats()

#     if bool(args.hostname):
#         get_hostnames()

#     if bool(args.verify):
#         verify_instance_setup()

#     if int(args.create):
#         create_instances(int(args.create))

#     if int(args.start):
#         start_instances(int(args.start))

#     if int(args.stop):
#         stop_instances(int(args.stop))

#     if int(args.terminate):
#         terminate_instances(int(args.terminate))

#     if int(args.setup):
#         setup_instances(int(args.setup))
