from flask import Flask, request, jsonify
app = Flask(__name__)


@app.route('/get', methods=["GET", "POST"])
def hello():
    if request.method == "POST":
        return jsonify({"Name": "Success"}), 200
    if request.method == "GET":
        return jsonify({"Name": "Success"}), 200
if __name__ == '__main__':
    app.run()