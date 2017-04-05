import json


with open("data.text") as f:
    data = json.load(f)

c = 0
t = {}

for k ,v in data.iteritems():
    if v["ReqBy"]:
        c += 1
        t[k] = v
    
print c
print t