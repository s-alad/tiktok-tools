import base64
import json
from flask import jsonify
import os
import io 
from PIL import Image



def get_image_bytes(fp:str):
    print(fp)
    img = Image.open(fp, mode='r')
    #print(img)
    imgByteArr = io.BytesIO()
    img.save(imgByteArr, format='JPEG')
    encoded = base64.encodebytes(imgByteArr.getvalue()).decode('ascii')
    #print(imgByteArr)
    #print(encoded)
    return encoded

def get_images():
    downloads_folder = "\downloads"
    cwd = os.getcwd()
    pth =cwd + downloads_folder
    encoded_images = []
    limit = 10
    total = 0
    for directory in os.listdir(pth):
        #check if slides folder exists in directory
        if os.path.exists(pth + "\\" + directory + "\\slides"):
            #get all images in slides folder
            for image in os.listdir(pth + "\\" + directory + "\\slides"):

                # check if file is video
                if image.endswith(".mp4"):
                    continue

                #get image bytes
                image_bytes = get_image_bytes(pth + "\\" + directory + "\\slides\\" + image)
                #append to list
                encoded_images.append(image_bytes)
                #check if limit is reached
        print(total)
        #if total == limit:    
        #    break
        #total += 1
    
    return {'dump': encoded_images}

#write out to json file
with open('data.json', 'w') as outfile:
    json.dump(get_images(), outfile)