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

instance_ids = []

ec2_client = boto3.client('ec2')
response = ec2_client.describe_instances()
for reservation in response["Reservations"]:
    for instance in reservation["Instances"]:
            print(instance["InstanceId"])
            instance_ids.append(instance["InstanceId"])

instance_ids = instance_ids[:60]
print(len(instance_ids))

