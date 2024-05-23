"""
Usage:
$ python3 calculate_verifier.py > calculate_verifier.csv
"""

from gen import rand_ascii, rand_hex

for _ in range(100):
    row = [
        rand_ascii(16), # username
        rand_ascii(16), # password
        rand_hex(32), # salt (32 bytes, little endian)
        "REPLACE_ME_IN_CSV", # expected value (32 bytes, little endian)
    ]
    print(",".join(row))
