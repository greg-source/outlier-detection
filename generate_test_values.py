import json
import random

num_entries = 200000
data = []
names = ["machine" + str(i) for i in range(1, 151)]

for i in range(num_entries):
    name = random.choice(names)
    age = random.choice(["%s years" % random.randint(1, 60),
                         "%s months" % random.randint(6, 360),
                         "%s days" % random.randint(90, 1500)])
    entry = {"name": name, "age": age}
    data.append(entry)

with open('data.json', 'w') as f:
    json.dump(data, f, indent=2)