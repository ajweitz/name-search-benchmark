
import random
import string

def search_task(client,path):
    for i in range(8):
        letters = string.ascii_lowercase
        input = ''.join(random.choice(letters) for _ in range(i))
        client.get(f"{path}?search={input}", name=path)