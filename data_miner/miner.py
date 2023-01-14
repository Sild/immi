import sys
import os
import json
import requests

def get_script_dir():
    return os.path.dirname(os.path.realpath(__file__))

def get_all_companies():
    companylist_path = "{}/../data/companylist.csv".format(get_script_dir())
    companylist = set()
    with open(companylist_path) as clp:
        clp.readline()
        line = clp.readline()
        while line:
            line_arr = line.split(",")
            companylist.add(line.split(",")[0].strip('"'))
            line = clp.readline()
    return sorted(companylist)

companies = get_all_companies()
cookies = {
    "B": "b7vsujdgcb8p1&b=3&s=k5",
    "GUC": "AQEBAQFgxvRgz0Ii7gTL",
    "A1": "d=AQABBCGjxWACEIaBXrYBU-OZMRqH3Zue_7MFEgEBAQH0xmDPYAAAAAAA_eMAAAcIIaPFYJue_7M&S=AQAAAii9BcIMKs0spHDaoVC0EVc",
    "A3": "d=AQABBCGjxWACEIaBXrYBU-OZMRqH3Zue_7MFEgEBAQH0xmDPYAAAAAAA_eMAAAcIIaPFYJue_7M&S=AQAAAii9BcIMKs0spHDaoVC0EVc",
    "A1S": "d=AQABBCGjxWACEIaBXrYBU-OZMRqH3Zue_7MFEgEBAQH0xmDPYAAAAAAA_eMAAAcIIaPFYJue_7M&S=AQAAAii9BcIMKs0spHDaoVC0EVc&j=WORLD"
}
#  Saturday, July 20, 1985 8:00:00 PM - Saturday, July 20, 2019 9:00:00 PM
url_template = "https://query1.finance.yahoo.com/v7/finance/download/{ticker}?period1=315532800&period2=1563656400&interval=1d&events=history&crumb=0BwI6UkuB1w"

total_count = len(companies)
current = 0
active = False
for ticker in companies:
    # if ticker == "ORIT":
    #     active = True
    # if not active:
    #     current += 1
    #     continue
    url = url_template.format(ticker=ticker)
    response = requests.get(url, cookies=cookies, headers={})
    response_data = response.text
    print(response_data)
    break 
    ticker_dir = "{}/../data/{}".format(get_script_dir(), ticker)
    if not os.path.exists(ticker_dir):
        os.makedirs(ticker_dir)
    ticker_data_path = "{}/history_data.csv".format(ticker_dir)
    with open(ticker_data_path, 'w') as output:
        output.write(response_data)
    current += 1
    print("{}/{}: {} history data ready: {}".format(current, total_count, ticker, ticker_data_path))
