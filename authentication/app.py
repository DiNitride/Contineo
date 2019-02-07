import json
from flask import Flask
app = Flask("Contine Mock Auth")

@app.route("/")
def hello():
    return "Hello World!"

@app.route("/authenticate")
def auth():
    return json.dumps({"session": "12345"})