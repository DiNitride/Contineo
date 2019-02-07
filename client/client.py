import websockets
import asyncio
import random
import string
import json
import signal

def signal_handler(sig, frame):
    print("exiting")
    exit(0)

signal.signal(signal.SIGINT, signal_handler)

async def main():
    ws = await websockets.connect("ws://localhost:3000/controller", extra_headers=[("token", "dinitride"), ("conn_type", "client")])
    
    accessToken = ws.response_headers["access_token"]
    print("access token " + accessToken)

    # Send init data
    await ws.send(json.dumps({
        "origin": "client",
        "type": "init",
        "init_token": accessToken
    }))

    # Wait for response
    print("waiting for resp")
    resp = await ws.recv()
    print(resp)

    # Send updates
    while True:
        resp = await ws.recv()
        print(resp)

asyncio.get_event_loop().run_until_complete(main())
