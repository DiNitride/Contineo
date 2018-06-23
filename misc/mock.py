import websockets
import asyncio
import random
import string
import json

key = "5eda2e3afed846d0905cd4c180919b75"

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
    ws = await websockets.connect("ws://localhost:3000/controller")

    # Send init data
    await ws.send(json.dumps({
        "origin": "worker",
        "type": "init",
        "key": key,
        "data": data
    }))

    # Wait for response
    resp = await ws.read()
    print(resp)

    # Send updates
    while True:
        refresh()
        await ws.send(json.dumps({"origin": "worker", "type": "update", "data": data}))
        print(f"sent {data}")
        await asyncio.sleep(5)

asyncio.get_event_loop().run_until_complete(main())
