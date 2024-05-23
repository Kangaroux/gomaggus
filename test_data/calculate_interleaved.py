"""
Usage:
$ python3 calculate_interleaved.py > calculate_interleaved.csv
"""

from gen import random_string

for _ in range(100):
    row = [
        random_string("hex", 64).upper(), # S key (32 bytes, little endian)
        "REPLACE_ME_IN_CSV", # expected value (40 bytes, little endian)
    ]
    print(",".join(row))
