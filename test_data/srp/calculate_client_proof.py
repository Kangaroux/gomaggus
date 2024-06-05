"""
Usage:
$ python3 calculate_client_proof.py > calculate_client_proof.csv
"""

from gen import rand_ascii, rand_hex

for _ in range(100):
    row = [
        rand_ascii(16).upper(), # username
        rand_hex(32), # salt
        rand_hex(32), # client public key
        rand_hex(32), # server public key
        rand_hex(40), # session key
        "REPLACE_ME_IN_CSV", # expected value
    ]
    print(",".join(row))
