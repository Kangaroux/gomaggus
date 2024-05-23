"""
Usage:
$ python3 calculate_x.py > calculate_x.csv
"""

from gen import random_string

for _ in range(100):
    row = [
        random_string("ascii", 16).upper(), # username
        random_string("ascii", 16).upper(), # password
        random_string("hex", 64).upper(), # salt (32 bytes)
        "EXPECTED_VALUE_HERE",
    ]
    print(",".join(row))
