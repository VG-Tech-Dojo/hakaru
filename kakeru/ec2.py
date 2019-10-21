import os

import boto3


REGION = os.environ.get('AWS_REGION', 'ap-northeast-1')


def list_autoscaling_instances(asg, self_addr):

    asclient = boto3.client('autoscaling', REGION)
    ec2client = boto3.client('ec2', REGION)

    return _list_autoscaling_instances(asg, asclient, ec2client, self_addr)


def _list_autoscaling_instances(asg, asclient, ec2client, self_addr):

    m = asclient.describe_auto_scaling_groups(AutoScalingGroupNames=[asg])

    instance_ids = [
        i['InstanceId']
        for g in m['AutoScalingGroups']
        for i in g['Instances']
        if g['AutoScalingGroupName'] == asg
    ]

    result = ec2client.describe_instances(InstanceIds=instance_ids)

    return {
        i['InstanceId']:dict(
            id=i['InstanceId'],
            private_name=i['PrivateDnsName'],
            private_addr=i['PrivateIpAddress'],
            public_name=i['PublicDnsName'],
            public_addr=i['PublicIpAddress']
        )
        for reservation in result['Reservations']
        for i in reservation['Instances']
        if self_addr != i['PrivateIpAddress']
    }


if __name__ == '__main__':

    import pprint

    sess = boto3.Session(profile_name='sunrise2018')
    autoscaling = sess.client('autoscaling', REGION)
    ec2 = sess.client('ec2', REGION)

    pprint.pprint(_list_autoscaling_instances('kakeru', autoscaling, ec2))
