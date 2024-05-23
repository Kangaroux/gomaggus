"""
Usage:
$ python3 calculate_u.py > calculate_u.csv
"""

from gen import random_string

for _ in range(100):
    row = [
        random_string("hex", 64).upper(), # client public key (32 bytes, little endian)
        random_string("hex", 64).upper(), # server public key (32 bytes, little endian)
        "REPLACE_ME_IN_CSV", # expected value (20 bytes, little endian)
    ]
    print(",".join(row))
