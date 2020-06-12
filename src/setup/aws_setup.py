import boto3
import os
import time
import argparse
import pprint 
import multiprocessing
import paramiko

parser = argparse.ArgumentParser(description="This program autocreates, starts, and stop Amazon EC2 Instancess.")
parser.add_argument('-create', default=0, help="Number of new instances to create")
parser.add_argument('-start', default=0, help="Number of new instances to start")
parser.add_argument('-stop', default=0, help="Number of new instances to stop")
parser.add_argument('-terminate', default=0, help="Number of new instances to terminate")
parser.add_argument('-stats', default=False, help="Number of new instances to terminate")
parser.add_argument('-run', default=True, help="Command to run on running instances")
parser.add_argument('-setup', default=True, help="Setup instances")

ec2_resource = boto3.resource('ec2',
    aws_access_key_id=os.environ['AWS_ACCESS_KEY_ID'],
    aws_secret_access_key=os.environ['AWS_SECRET_ACCESS_KEY'],
    region_name='us-west-1')

ec2_client = boto3.client('ec2',
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
    # for line in stdout.read():
    #     print(line)
    # if stdout.channel.recv_exit_status():
    #     print(ip, " failed setup")
    
    client.close()

def setup_instances(num):
    running_instances = list(ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}]))
    if len(running_instances) < num:
        print("Not all instances are up yet!")
        return

    running_instance_ips = [instance.public_dns_name for instance in running_instances]
    pool = multiprocessing.Pool(len(running_instance_ips))
    pool.map(worker, running_instance_ips)
    
    # for instance in running_instances:
    #     client = paramiko.SSHClient()
    #     client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    #     client.connect(hostname=instance.public_dns_name, username="ubuntu", pkey=key)
    #     stdin, stdout, stderr = client.exec_command('git clone https://github.com/edwjchen/agario.git')
    #     stdin.flush()
    #     if stdout.channel.recv_exit_status():
    #         print(instance.public_ip_address, " failed clone")

    #     stdin, stdout, stderr = client.exec_command('cd agario/src/setup')
    #     stdin.flush()
    #     if stdout.channel.recv_exit_status():
    #         print(instance.public_ip_address, " failed cd?")

    #     stdin, stdout, stderr = client.exec_command('. setup.sh')
    #     stdin.flush()
    #     if stdout.channel.recv_exit_status():
    #         print(instance.public_ip_address, " failed setup?")
        
    pass
    
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
    args = parser.parse_args()
    if int(args.create):
        create_instances(int(args.create))

    if int(args.start):
        start_instances(int(args.start))

    if int(args.stop):
        stop_instances(int(args.stop))

    if int(args.terminate):
        terminate_instances(int(args.terminate))

    if int(args.setup):
        setup_instances(int(args.setup))

    if bool(args.stats):
        get_stats()

    



