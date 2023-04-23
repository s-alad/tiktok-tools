import requests
import json

class Tik:

    def __init__(self) -> None:
        pass

    def tiklydown(vidurl: str):
        r = requests.get("https://api.tiklydown.me/api/download?url={}".format(vidurl))
        if r.status_code == 200:
            return r.json()

    def no_watermark(data):
        return data["video"]["noWatermark"]
    
    def create(url, path):
        r = requests.get(url)
        with open(path, "wb") as f:
            f.write(r.content)
    
    def download(vidurl: str, path: str):
        u = Tik.tiklydown(vidurl)
        Tik.create(Tik.no_watermark(u), path)

if __name__ == "__main__":
    u = Tik.download("https://www.tiktok.com/t/ZTRTfYL5w/", "test.mp4")