import requests
import random
from joblib import Parallel, delayed


url = 'http://localhost:8081/hakaru'

def sendRequest():
    payload = {'name': 'pykakeru', 'value': random.randint(0, 2000000)}
    response = requests.get(url, params=payload)
    print(response.status_code)

    return response

if __name__ == "__main__":
    Parallel(n_jobs=10)([delayed(sendRequest)() for i in range(100000)])
