import json


with open("data.text") as f:
    data = json.load(f)

print len(data)
