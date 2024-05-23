"""
Usage:
$ python3 calculate_client_s.py > calculate_client_s.csv
"""

from gen import random_string

for _ in range(100):
    row = [
        random_string("hex", 64).upper(), # server public key (32 bytes, little endian)
        random_string("hex", 64).upper(), # client private key (32 bytes, little endian)
        random_string("hex", 64).upper(), # x (32 bytes, little endian)
        random_string("hex", 64).upper(), # u (32 bytes, little endian)
        "REPLACE_ME_IN_CSV", # expected value (32 bytes, little endian)
    ]
    print(",".join(row))
