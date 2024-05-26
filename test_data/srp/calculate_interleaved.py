"""
Usage:
$ python3 calculate_interleaved.py > calculate_interleaved.csv
"""

from gen import rand_hex

for _ in range(100):
    row = [
        rand_hex(32), # S key (32 bytes, little endian)
        "REPLACE_ME_IN_CSV", # expected value (40 bytes, little endian)
    ]
    print(",".join(row))
