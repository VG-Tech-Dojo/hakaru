import json

import slave

REPORT_URL = 'http://sunrise2018-kakeru-report.s3-website-ap-northeast-1.amazonaws.com/'


def param(event, key, default=None):

    return event.get('queryStringParameters', {}).get(key, default)


def response(status, body):

    return dict(
        statusCode=status,
        body=json.dumps(body)
    )


def execute(host, port, scenario, slaves):

    return slave.execute([
        'make',
        '-C',
        '/opt/sunrise2018/kakeru',
        'ecs/multinode',
        'HOST='+host,
        'PORT='+str(port),
        'SCENARIO='+str(scenario),
        'SLAVES='+str(slaves)
    ])


def handler(event, context):

    host = param(event, 'host')
    port = param(event, 'port', '80')
    slaves = param(event, 'slaves', '4')
    scenario = param(event, 'scenario', '1')

    if not host:
        return response(400, dict(error='invalid host: {0}'.format(host)))

    try:
        port = int(port)
    except ValueError:
        return response(400, dict(error='invalid port: {0}'.format(port)))

    try:
        slaves = int(slaves)
    except ValueError:
        return response(400, dict(error='invalid number format: {0}'.format(slaves)))

    try:
        scenario = int(scenario)
    except ValueError:
        return response(400, dict(error='invalid number format: {0}'.format(scenario)))

    execute(host, port, scenario, slaves)

    report_url = '{root}{host}:{port}/index.html'.format(root=REPORT_URL, host=host, port=port)

    return response(
        200,
        dict(
            msg='launched',
            report_url=report_url
        ))
