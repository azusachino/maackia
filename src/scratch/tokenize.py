filepath = "assets/the_verdict.txt"

with open(filepath, "r", encoding="utf-8") as f:
    raw_text = f.read()
print("total number of the character: ", len(raw_text))
print(raw_text[:99])
