import sys
import os
import json
from urllib import request

if len(sys.argv) < 2:
    print("output dir must be specified as $1")
    exit(1)
OUTPUT_DIR = sys.argv[1]

def get_all_companies():
    script_dir = os.path.dirname(os.path.realpath(__file__))
    companylist_path = "{}/../data/companylist.csv".format(script_dir)
    companylist = []
    with open(companylist_path) as clp:
        clp.readline()
        line = clp.readline()
        while line:
            line_arr = line.split(",")
            companylist.append(line.split(",")[0].strip('"'))
            line = clp.readline()
    return companylist

companies = get_all_companies()
print(companies)
