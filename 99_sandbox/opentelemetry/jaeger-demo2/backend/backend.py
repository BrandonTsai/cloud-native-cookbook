import logging

from flask import Flask, jsonify, request

app = Flask(__name__)

USERS = ["Brandon", "Felix", "Adrian", "Neil"]


@app.route('/user', methods=['GET'])
def get_user():
    uid = request.json.get('id')
    logging.getLogger().debug("Get user id `%s` from frontend" % uid)

    resp = {"name": USERS[int(uid)]}
    return jsonify(resp), 200


if __name__ == '__main__':
    app.run(host='127.0.0.1', port=8081)
