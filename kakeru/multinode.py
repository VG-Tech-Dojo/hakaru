import argparse
import glob
import json
import os
import subprocess
import sys
import tempfile
import time
import traceback
import urllib.request

import boto3

import ec2
import slave

webhook = 'incoming webhook url'
METADATA = 'http://169.254.169.254/latest/meta-data/local-ipv4'


def send_slack(msg, host, scenario, slaves, task_arn):

    body = dict(text='''
{msg}
host: {host}
slaves: {slaves}
scenario: {scenario}
task_arn: {task_arn}
    '''.format(**locals()))

    headers = {
        'Content-Type': 'application/json; charset=utf-8',
    }

    urllib.request.urlopen(webhook, json.dumps(body).encode('utf-8'))


def make_parser():

    p = argparse.ArgumentParser()

    add = p.add_argument

    add('--host', type=str, required=True)
    add('--port', type=int, default=80)
    add('--scenario', type=str, default=1)
    add('--logdir', type=str, default='../log')
    add('--bucket', type=str, default='sunrise2018-kakeru-report')
    add('--asg', type=str, default='kakeru')
    return p


TEMPLATE_PATH = 'scenarios/hakaru-test-multinode-{scenario}.xml.template'
OUTPUT_PATH = 'tmp/kakeru-test.xml'

CLIENT_TEMPLATE = '''<client host="{addr}" maxusers="40000" />
'''


def make_test_xml(instances, host, port, scenario):

    addrs = [x['private_name'] for x in instances.values() if 'private_name' in x]

    print(addrs)

    clients = ''.join([CLIENT_TEMPLATE.format(addr=addr) for addr in addrs])

    with open(TEMPLATE_PATH.format(scenario=scenario), 'r', encoding='utf-8') as fp:
        data = fp.read()
        data = data.format(clients=clients, host=host, port=port)
        print(data.encode('utf-8'))

    if not os.path.exists('tmp'):
        os.makedirs('tmp')

    with open(OUTPUT_PATH, 'w', encoding='utf-8') as fp:
        fp.write(data)

    return OUTPUT_PATH


def self_ipaddr():

    return urllib.request.urlopen(METADATA).read().decode('utf-8')


def task_arn():

    metadata = urllib.request.urlopen(METADATA).read()
    metadata = json.loads(metadata)

    arn = metadata.get('TaskARN')

    return arn


def tsung(xml, logdir):

    addr = self_ipaddr()

    p = subprocess.Popen(['tsung', '-I', addr, '-l', logdir, '-f', xml, 'start'])
    p.wait()

    for d in glob.glob(logdir+'/*'):
        dp = subprocess.Popen(['/usr/lib/x86_64-linux-gnu/tsung/bin/tsung_stats.pl'], cwd=d)
        dp.wait()


def upload(host, port, logdir, bucket):

    path = 's3://{bucket}/{host}:{port}/'.format(**locals())

    p = subprocess.Popen(['aws', 's3', 'cp', '--recursive', logdir, path])

    return p.wait()


def make_index(host, port, bucket):

    s3 = boto3.client('s3')
    response = s3.list_objects_v2(Bucket=bucket, Prefix='{host}:{port}/'.format(host=host, port=port), Delimiter='/')

    print(response)

    prefixes = [p['Prefix'] for p in response.get('CommonPrefixes', [])]

    with tempfile.NamedTemporaryFile() as tmp:

        tmp.write('''
<html><head><title>Reports for {host}:{port}</title></head>
<h1>Reports for {host}:{port}</h1>
<ul>
        '''.format(**locals()).encode('utf-8'))

        for prefix in sorted(prefixes):
            print(prefix)
            dirname = prefix.split('/')[1]
            tmp.write('''
<li>
  <a href="/{prefix}report.html">
    {dirname}
  </a>
</li>
'''.format(**locals()).encode('utf-8'))

        tmp.write(b'''
</ul>
''')

        tmp.flush()
        tmp.seek(0)

        key = '{host}:{port}/index.html'.format(**locals())
        s3.put_object(Bucket=bucket, Key=key, Body=tmp, ContentType='text/html; charset=utf-8')

        print('create index:', key)


def main():

    arg = make_parser().parse_args()
    addr = self_ipaddr()

    print(arg)

    print('準備中')

    instances = ec2.list_autoscaling_instances(arg.asg, addr)
    slave_addrs = [x['private_addr'] for x in instances.values()]

    try:
        slave.keygen(slave_addrs)

        xmlpath = make_test_xml(instances, arg.host, arg.port, arg.scenario)

        print('開始')
        tsung(xmlpath, arg.logdir)

        print('レポート作成中')

        upload(arg.host, arg.port, arg.logdir, arg.bucket)

        make_index(arg.host, arg.port, arg.bucket)

        print('おしまい')

    except Exception as e:
        traceback.print_exc()
        print('エラー出た')


if __name__ == '__main__':
    main()
