import sys
import uuid
import json

name = sys.argv[1]

game = input("Game: ")
gamemode = input("Gamemode: ")
desc = input("Description: ")
id = uuid.uuid4().hex
key = uuid.uuid4().hex

template = {
    "game": game,
    "gamemode": gamemode,
    "desc": desc,
    "key": key,
    "id": id
}

with open(f"workers\{name}.json", "w") as file:
    file.write(json.dumps(template))
