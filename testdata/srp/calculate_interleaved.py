"""
Usage:
$ python3 calculate_interleaved.py > calculate_interleaved.csv
"""

from gen import rand_hex

for _ in range(100):
    row = [
        rand_hex(32), # S key
        "REPLACE_ME_IN_CSV", # expected value
    ]
    print(",".join(row))
