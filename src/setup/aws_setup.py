import boto3
import os
import time
import argparse
import pprint
import multiprocessing
import paramiko
import sys

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
EXPERIMENT_NAME= ""

SERVERS_RUNNING = [] # list of dns names

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
        print(ip, " failed clone")
    client.close()


def refresh_instances():
    running_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]))
    running_instance_ips = [instance.public_dns_name for instance in running_instances]
    pool = multiprocessing.Pool(len(running_instance_ips))
    pool.map(refresh, running_instance_ips)

def start_entry():
    global ENTRY_NAME
    running_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]))
    if not len(running_instances):
        print("!!Error!!: No instances currently running!")

    entry_instance = running_instances[0]
    ENTRY_NAME = entry_instance.public_dns_name

    key = paramiko.RSAKey.from_private_key_file("the-key-to-her-heart.pem")
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    print("Starting entryserver on", ENTRY_NAME)
    client.connect(hostname=ENTRY_NAME, username="ubuntu", pkey=key)
    stdin, stdout, stderr = client.exec_command('export GOPATH=/home/ubuntu/agario; /usr/local/go/bin/go run agario/src/peer_to_peer/entryserver.go')
    stdin.flush()
    # if stdout.channel.recv_exit_status():
        # print(ENTRY_NAME, " failed to start entryserver")
    time.sleep(2)
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
            print("data refresh")
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
        elif cmd_type == "wait":
            time.sleep(int(cmds[1]))
        elif cmd_type == "refresh":
            refresh_instances()
        elif cmd_type == "start_entry":
            start_entry()
        elif cmd_type == "kill_entry":
            print(ENTRY_NAME)
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
