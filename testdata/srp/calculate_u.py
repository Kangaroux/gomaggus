"""
Usage:
$ python3 calculate_u.py > calculate_u.csv
"""

from gen import rand_hex

for _ in range(100):
    row = [
        rand_hex(32), # client public key
        rand_hex(32), # server public key
        "REPLACE_ME_IN_CSV", # expected value
    ]
    print(",".join(row))
