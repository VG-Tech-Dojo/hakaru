import os
import subprocess
import time

import boto3

REGION = 'ap-northeast-1'
CLUSTER = 'arn:aws:ecs:ap-northeast-1:139990579284:cluster/kakeru-san'
SUBNETS = ['subnet-01b4f71d26c32b559', 'subnet-0c8bd552a7413a629']
SEURITY_GROUPS = ['sg-019dd505eb986a1df']

TASK_DEF = 'kakeru-fargate'

ecs = boto3.client('ecs', REGION)


def check(r):

    if r['failures']:
        raise Exception(r['failures'])


class Container(object):

    def __init__(self, data):
        self.data = data

    @property
    def private_addr(self):

        for i in self.data['networkInterfaces']:

            if 'privateIpv4Address' in i:
                return i['privateIpv4Address']


class Task(object):

    def __init__(self, ecs, data):
        self.ecs = ecs
        self.data = data

    @property
    def status(self):
        return self.last_status

    @property
    def is_running(self):
        return self.last_status == 'RUNNING'

    @property
    def last_status(self):
        return self.data['lastStatus']

    @property
    def arn(self):
        return self.data['taskArn']

    @property
    def cluster(self):
        return self.data['clusterArn']

    def update(self):

        result = self.ecs.describe_tasks(
            cluster=self.cluster,
            tasks=[self.arn])

        check(result)

        self.data = result['tasks'][0]

    def kill(self, reason='No reason'):

        result = self.ecs.stop_task(
            cluster=self.cluster,
            task=self.arn,
            reason=reason
        )

        print('kill:', self.arn)

        return result

    def wait_for_running(self):

        while not self.is_running:
            self.update()
            time.sleep(1)
            print('wait:', self.arn)

    @property
    def containers(self):

        return [Container(c) for c in self.data.get('containers', [])]

    @property
    def private_addr(self):

        for c in self.containers:

            if c.private_addr:
                return c.private_addr


def execute(cmdline, num=1):

    ret = []

    try:
        for n in range(num):
            result = ecs.run_task(
                cluster=CLUSTER,
                taskDefinition=TASK_DEF,
                launchType='FARGATE',
                networkConfiguration=dict(
                    awsvpcConfiguration=dict(
                        subnets=SUBNETS,
                        securityGroups=SEURITY_GROUPS,
                        assignPublicIp='ENABLED'
                    )
                ),
                overrides=dict(
                    containerOverrides=[
                        dict(
                            name='kakeru-san',
                            command=cmdline,
                        )
                    ]
                )
            )

            check(result)

            ret.extend([Task(ecs, x) for x in result['tasks']])
    except Exception as e:
        import traceback
        traceback.print_exc()
        killall(ret)

    return ret


class EC2Instance(object):

    def __init__(self, data):

        self.data = data

    @property
    def private_addr(self):
        return self.data.get('private_addr')


def waitall(tasks):

    for task in tasks:
        task.wait_for_running()

def killall(tasks):

    for t in tasks:
        t.kill()


def _keygen(tasks):

    script = os.path.join(os.path.dirname(__file__), 'keygen.sh')

    for task in tasks:
        subprocess.Popen(['/bin/sh', script, task.private_addr])


def keygen(addrs):

    script = os.path.join(os.path.dirname(__file__), 'keygen.sh')

    processes = [
        subprocess.Popen(['/bin/sh', script, private_addr])
        for private_addr in addrs]

    for p in processes:
        p.wait()
