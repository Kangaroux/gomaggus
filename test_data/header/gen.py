import random
import string

ascii_chars = string.ascii_letters + string.digits
hex_chars = "0123456789ABCDEF"

def rand_hex(slen: int) -> str:
    return "".join(random.choices(hex_chars, k=slen*2))

def rand_ascii(slen: int) -> str:
    return "".join(random.choices(ascii_chars, k=slen))
