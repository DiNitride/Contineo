from random import choice
from string import ascii_lowercase
import json
from flask import Flask, request
app = Flask("Contine Mock Auth")

allowed = ["dinitride"]

def gen_tok():
    return "tok_" + "".join([choice(ascii_lowercase) for i in range(10)])

@app.route("/")
def hello():
    return "Hello World!"

@app.route("/authenticate", methods=["POST"])
def auth():
    given_token = request.form["token"]
    print("given token " + given_token)
    if given_token in allowed:
        return json.dumps({"access_token": gen_tok()})
    else:
        return json.dumps({"access_token": "unauthorised"})