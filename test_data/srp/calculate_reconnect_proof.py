"""
Usage:
$ python3 calculate_reconnect_proof.py > calculate_reconnect_proof.csv
"""

from gen import rand_ascii, rand_hex

for _ in range(100):
    row = [
        rand_ascii(16).upper(), # username
        rand_hex(16), # random client data
        rand_hex(16), # random server data
        rand_hex(40), # session key
        "REPLACE_ME_IN_CSV", # expected value
    ]
    print(",".join(row))
