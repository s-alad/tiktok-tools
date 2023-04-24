import base64
import json
from flask import Flask
from flask import jsonify
import os
import io 
from PIL import Image
from flask_cors import CORS, cross_origin


app = Flask(__name__)
cors = CORS(app)
app.config['CORS_HEADERS'] = 'Content-Type'


@app.route("/")
def index():
    return '/'

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

@app.route("/getimagesquick", methods=['GET'])
@cross_origin()
def get_images_quick():
    #load data.json 
    with open('data.json') as json_file:
        data = json.load(json_file)
        return jsonify(data)


""" @app.route("/getimages", methods=['GET'])
@cross_origin()
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
    
    return jsonify({'dump': encoded_images})
 """
#all all cors

if __name__ == "__main__":
    app.run(debug=True)