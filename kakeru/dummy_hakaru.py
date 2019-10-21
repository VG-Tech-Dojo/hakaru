def hakaru(environ, start_response):

    start_response('200 OK',
                  [('Content-Type', 'text/plain; charset=utf-8')])

    return ['hello']


if __name__ == '__main__':

    from wsgiref import simple_server
    simple_server.make_server('', 8888, hakaru).serve_forever()
