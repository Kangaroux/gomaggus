"""
Usage:
$ python3 calculate_server_s.py > calculate_server_s.csv
"""

from gen import rand_hex

for _ in range(100):
    row = [
        rand_hex(32), # client public key (32 bytes, little endian)
        rand_hex(32), # server private key (32 bytes, little endian)
        rand_hex(32), # verifier (32 bytes, little endian)
        rand_hex(32), # u (32 bytes, little endian)
        "REPLACE_ME_IN_CSV", # expected value (32 bytes, little endian)
    ]
    print(",".join(row))
