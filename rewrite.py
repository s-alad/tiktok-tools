import requests
import json

def rewrite(link: str):
    cookies = {
        '',
    }

    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/111.0',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Accept-Encoding': 'gzip, deflate, br',
        'DNT': '1',
        'Upgrade-Insecure-Requests': '1',
        'Sec-Fetch-Dest': 'document',
        'Sec-Fetch-Mode': 'navigate',
        'Sec-Fetch-Site': 'none',
        'Sec-Fetch-User': '?1',
        'Sec-GPC': '1',
        'Connection': 'keep-alive',
        # 'TE': 'trailers',
    }

    full = 'https://www.tiktok.com/oembed?url={}'.format(link)
    print(full)

    oembed = requests.get(
        full,
        cookies=cookies,
        headers=headers,
    )
    response = requests.get(oembed)
    status = response.status_code
    #data = response.json()
    print(status)

    #print(data['author_url'], data["embed_product_id"])

    #pretty print data
    #print(json.dumps(data, indent=4, sort_keys=True))

#open up raw.txt
with open("raw.txt", "r") as f:
    for line in f:
        print(line)
        rewrite(line.strip())
