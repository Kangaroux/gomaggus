"""
Usage:
$ python3 calculate_server_proof.py > calculate_server_proof.csv
"""

from gen import rand_hex

for _ in range(100):
    row = [
        rand_hex(32), # client public key
        rand_hex(20), # client proof
        rand_hex(40), # session key
        "REPLACE_ME_IN_CSV", # expected value
    ]
    print(",".join(row))
