import websockets
import asyncio
import random
import string
import json

key = "worker_test"

data = {
    "latency": 0,
    "players": 0,
    "online": False
}


def refresh():
    data["online"] = random.choice([True, False])
    data["players"] = random.randint(0, 100)
    data["latency"] = random.randint(20, 50)


async def main():
    refresh()
    ws = await websockets.connect("ws://localhost:3000/controller", extra_headers=[("conn_type", "worker")])

    # Send init data
    await ws.send(json.dumps({
        "origin": "worker",
        "type": "init",
        "init_token": "worker_test",
        "data": data
    }))

    # Wait for response
    print("waiting for resp")
    resp = await ws.recv()
    print(resp)

    # Send updates
    while True:
        refresh()
        await ws.send(json.dumps({"origin": "worker", "type": "update", "data": data}))
        print(f"sent {data}")
        await asyncio.sleep(5)

asyncio.get_event_loop().run_until_complete(main())
