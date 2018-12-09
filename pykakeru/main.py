import requests
import random


url = 'http://localhost:8081/hakaru'

def sendRequest():
    payload = {'name': 'pykakeru', 'value': random.uniform(0, 200000)}
    response = requests.get(url, params=payload)
    print(response.status_code)

    return response

if __name__ == "__main__":
    for n in range(0, 10):
        sendRequest()
