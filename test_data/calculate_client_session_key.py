"""
Usage:
$ python3 calculate_client_session_key.py > calculate_client_session_key.csv
"""

from gen import rand_ascii, rand_hex

for _ in range(100):
    row = [
        rand_ascii(16), # username
        rand_ascii(16), # password
        rand_hex(32), # client public key (32 bytes, little endian)
        rand_hex(32), # client private key (32 bytes, little endian)
        rand_hex(32), # server public key (32 bytes, little endian)
        rand_hex(32), # salt (32 bytes, little endian)
        "REPLACE_ME_IN_CSV", # expected value (20 bytes, little endian)
    ]
    print(",".join(row))
