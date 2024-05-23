import random
import string
from typing import Literal

ascii = string.ascii_letters + string.digits
hex = string.hexdigits

def random_string(kind: Literal["ascii", "hex"], slen: int) -> str:
    return "".join(random.choices(ascii if kind == "ascii" else hex, k=slen))
