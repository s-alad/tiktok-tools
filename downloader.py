import asyncio
import io
import glob
import logging
import os
import time
import aiohttp
import urllib.request

from os import path
from multiprocessing import Process

from tiktokapipy.async_api import AsyncTikTokAPI
from tiktokapipy.models.video import Video
from tik import Tik

async def save_slideshow(video: Video, direct: str):

    # create the directory to store the slides and audio
    directory = direct+"\slides"
    if not path.exists(directory): os.mkdir(directory)

    # this filter makes sure the images are padded to all the same size
    vf = "\"scale=iw*min(1080/iw\,1920/ih):ih*min(1080/iw\,1920/ih)," \
         "pad=1080:1920:(1080-iw)/2:(1920-ih)/2," \
         "format=yuv420p\""

    # downloads the images
    for i, image_data in enumerate(video.image_post.images):
        url = image_data.image_url.url_list[-1]
        urllib.request.urlretrieve(url, path.join(directory, f"slide_{i:02}.jpg"))

    # downloads the audios
    urllib.request.urlretrieve(video.music.play_url, path.join(direct, f"audio.mp3"))

    # use ffmpeg to join the images and audio
    command = [
        "ffmpeg",
        "-r 2/5",
        f"-i {directory}/slide_%02d.jpg",
        f"-i {direct}/audio.mp3",
        "-r 30",
        f"-vf {vf}",
        "-acodec copy",
        f"-t {len(video.image_post.images) * 2.5}",
        f"{directory}/video.mp4",
        "-y"
    ]
    ffmpeg_proc = await asyncio.create_subprocess_shell(
        " ".join(command),
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )
    _, stderr = await ffmpeg_proc.communicate()
    generated_files = glob.glob(path.join(directory, f"*"))

    # check if the video was generated
    if not path.exists(path.join(directory, f"video.mp4")):
        print(path.join(directory, f"video.mp4"))
        # optional ffmpeg logging step
        logging.error(stderr.decode("utf-8"))
        for file in generated_files:
            os.remove(file)
        raise Exception("Something went wrong with piecing the slideshow together")

    # temporarily store the video and return it
    with open(path.join(directory, f"video.mp4"), "rb") as f:
        ret = io.BytesIO(f.read())
    
    # remove video.mp4 in slides folder
    os.remove(path.join(directory, f"video.mp4"))

    return ret

async def save_video(video: Video, direct: str, link: str):

    # attempt to download the audio
    try: urllib.request.urlretrieve(video.music.play_url, path.join(direct, f"audio.mp3"))
    except: print("No audio found")

    # try to download without watermark
    try:
        Tik.download(link, path.join(direct, f"video.mp4"))
        return False
    except:
        print("No video found")
        async with aiohttp.ClientSession() as session:
            async with session.get(video.video.download_addr) as resp:
                return io.BytesIO(await resp.read())

async def download_video(link:str):

    quick = link.split("/")
    quick_author = quick[3].replace("@", "")
    quick_id = quick[5].split("?")[0]
    quick_direct = "downloads" + "\\" + str(quick_author) + "_" + str(quick_id)

    if  path.exists(path.join(quick_direct, f"video.mp4")) and \
        path.exists(path.join(quick_direct, f"desc.txt"))  and \
        path.exists(path.join(quick_direct, f"link.txt")):
        print(" >> already exists")
        return

    # mobile emulation is necessary to retrieve slideshows
    # if you don't want this, you can set emulate_mobile=False and skip if the video has an image_post property
    async with AsyncTikTokAPI(emulate_mobile=True, navigation_timeout=120) as api:

        print("=" * 100)
        video: Video = await api.video(link)
        author = video.author
        desc = video.desc
        print(f" >> downloading video {video.id} by {author}")

        direct = "downloads" + "\\" + str(author.unique_id) + "_" + str(video.id)
        if not path.exists(direct): os.mkdir(direct)

        #check if video.mp4, desc.txt, and link.txt already exist in the direct folder

        if video.image_post:
            print(" >> type slideshow")
            downloaded = await save_slideshow(video, direct)
        else:
            print(" >> type video")
            downloaded = await save_video(video, direct, link)

        if (type(downloaded) != bool):
            with open(path.join(direct, f"video.mp4"), "wb") as f:
                f.write(downloaded.read())
        
        with open(path.join(direct, f"desc.txt"), "w", encoding="utf-8") as f:
            try:
                f.write(desc)
            except: 
                print(" >> failed to write description")
                f.write("failed to write description: link: " + link)
        
        with open(path.join(direct, f"link.txt"), "w", encoding="utf-8") as f:
            f.write(link)
        
        print(" >> done")
        print("=" * 100)


#open links.txt and read the links
f = open("links/results.txt", "r")
links = f.readlines()
f.close()

#loop through the links and download the videos
start = time.time()
for link in links:
    #print(link)
    try:
        asyncio.run(download_video(link))
    except:
        print(" >> unresolved process failure")
        print("=" * 100)
end = time.time()
print("Total time taken: ", end - start)


""" loop = asyncio.get_event_loop()
loop.run_until_complete(asyncio.gather(*[download_video(link) for link in links]))
loop.close()
 """

""" from playwright.async_api import async_playwright

def run_download_video(link):
    asyncio.run(download_video(link))

if __name__ == "__main__":

    track_time_start = time.time()
    

    processes = []
    for link in links:
        p = Process(target=run_download_video, args=(link,))
        processes.append(p)

    for p in processes:
        p.start()
    
    for p in processes:
        p.join()
    
    track_time_end = time.time()
    print("Total time taken: ", track_time_end - track_time_start) """