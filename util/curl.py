import sys,os,time

#python3 util/curl.py  http://localhost:8081 10

dist = sys.argv[1]
n = int(sys.argv[2])

for i in range(n):
    os.system("curl -v \"%s/hakaru?name=\"hoge\"&value=50\"" % dist)
    time.sleep(0.2)
