import os
import requests
import logging
from flask import Flask, request

app = Flask(__name__)


@app.route('/')
def hello_world():

    uid = request.args.get('uid', default=0)
    logging.getLogger().debug("Get uid `%s` from user" % uid)
    backend_url = os.environ.get(
        'BACKEND_URL', default="http://localhost:8081")
    backend_api_user = f'{backend_url}/user'
    response = requests.get(backend_api_user, json={"id": uid})

    if response.status_code == 200:
        user_data = response.json()
        username = user_data.get('name')
        return f"Hello {username}", 200
    else:
        return "Backend Error", 500


if __name__ == '__main__':
    app.run(host='127.0.0.1', port=8080)
