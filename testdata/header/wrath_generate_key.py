"""
Usage:
$ python3 wrath_generate_key.py > wrath_generate_key.csv
"""

from gen import rand_hex

for _ in range(100):
    row = [
        rand_hex(40), # session key (40 bytes, little endian)
        rand_hex(16), # fixed key (40 bytes, little endian)
        "REPLACE_ME_IN_CSV", # expected value (20 bytes, little endian)
    ]
    print(",".join(row))
