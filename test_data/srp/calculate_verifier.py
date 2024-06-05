"""
Usage:
$ python3 calculate_verifier.py > calculate_verifier.csv
"""

from gen import rand_ascii, rand_hex

for _ in range(100):
    row = [
        rand_ascii(16).upper(), # username
        rand_ascii(16).upper(), # password
        rand_hex(32), # salt
        "REPLACE_ME_IN_CSV", # expected value
    ]
    print(",".join(row))
