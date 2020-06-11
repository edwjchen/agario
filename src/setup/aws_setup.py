import boto3
import os

ec2_resource = boto3.resource('ec2',
    aws_access_key_id=os.environ['AWS_ACCESS_KEY_ID'],
    aws_secret_access_key=os.environ['AWS_SECRET_ACCESS_KEY'],
    region_name='us-west-1')
    
#create instances
# instance = ec2_resource.create_instances(ImageId='ami-0318e6f2445586bd7', 
#     InstanceType='t2.micro', 
#     KeyName='the-key-to-her-heart', 
#     MinCount=1, 
#     MaxCount=60,
#     Monitoring={'Enabled': False},
#     SecurityGroups=['agario'])

# running_instances = ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['running']}])

# for instance in running_instances:
#     print(instance.id, instance.instance_type)
#     ec2_resource.instances.filter(InstanceIds=[instance.id]).stop()

stopped_instances = ec2_resource.instances.filter(Filters=[{'Name': 'instance-state-name', 'Values': ['stopped', 'stopping']}])
for instance in stopped_instances:
    print(instance.id, instance.instance_type)
    print(instance.private_ip_address)
    print()


# instance_ids = []

# ec2_client = boto3.client('ec2')
# response = ec2_client.describe_instances()
# for reservation in response["Reservations"]:
#     for instance in reservation["Instances"]:
#             print(instance["InstanceId"])
#             print(instance)
#             print()
#             instance_ids.append(instance["InstanceId"])
