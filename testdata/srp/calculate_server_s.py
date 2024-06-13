"""
Usage:
$ python3 calculate_server_s.py > calculate_server_s.csv
"""

from gen import rand_hex

for _ in range(100):
    row = [
        rand_hex(32), # client public key
        rand_hex(32), # server private key
        rand_hex(32), # verifier
        rand_hex(32), # u
        "REPLACE_ME_IN_CSV", # expected value
    ]
    print(",".join(row))
