import zlib, json, base64

ZIPJSON_KEY = 'base64(zip(o))'

def json_zip(j):

    j = {
        ZIPJSON_KEY: base64.b64encode(
            zlib.compress(
                json.dumps(j).encode('utf-8')
            )
        ).decode('ascii')
    }

    return j

#load data.json
with open('data.json') as json_file:
    data = json.load(json_file)

#compress data
compressed = json_zip(data)

#write out to json file
with open('cdata.json', 'w') as outfile:
    json.dump(compressed, outfile)