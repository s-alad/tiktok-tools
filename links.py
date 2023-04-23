import collections
import os
import pickle
import time
import requests

from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.common.exceptions import TimeoutException
from selenium.webdriver.chrome.options import Options as ChromeOptions
from selenium.webdriver.chrome.service import Service as ChromeService
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.wait import WebDriverWait
from webdriver_manager.chrome import ChromeDriverManager

collections.Callable = collections.abc.Callable

class Linker:
    path = os.getcwd()
    def __init__(self):
        options = ChromeOptions()
        # add experimental option to keep window open after test is done for debugging
        options.add_experimental_option("detach", True)
        prefs = {'profile.default_content_setting_values': {'cookies': 2, 'images': 2, 'javascript': 2, 
                'plugins': 2, 'popups': 2, 'geolocation': 2, 
                'notifications': 2, 'auto_select_certificate': 2, 'fullscreen': 2, 
                'mouselock': 2, 'mixed_script': 2, 'media_stream': 2, 
                'media_stream_mic': 2, 'media_stream_camera': 2, 'protocol_handlers': 2, 
                'ppapi_broker': 2, 'automatic_downloads': 2, 'midi_sysex': 2, 
                'push_messaging': 2, 'ssl_cert_decisions': 2, 'metro_switch_to_desktop': 2, 
                'protected_media_identifier': 2, 'app_banner': 2, 'site_engagement': 2, 
                'durable_storage': 2}}
        options.add_experimental_option('prefs', prefs)
        options.add_argument("disable-infobars")
        options.add_argument("--disable-extensions")
        #options.add_argument("--headless=new")
        options.add_argument("user-data-dir={}".format(self.path + "/profile"))
        self.driver = webdriver.Chrome(service=ChromeService(ChromeDriverManager().install()), options=options)

    def get_link(self, url):
        driver = self.driver
        driver.get(url)
        #wait until url changes to a different url from the original
        WebDriverWait(driver, 10).until(EC.url_changes(url))
        return driver.current_url
    
    def get_links(self, links):
        new_links = []
        for link in links:
            print(" >> getting link for {}".format(link))
            try:
                new_links.append(self.get_link(link))
            except:
                print(" >> failed to get link for")
        return new_links
    
    def fast_links(self, links):
        driver = self.driver
        for link in links:
            print(" >> opening link {}".format(link))
            time.sleep(.2)
            driver.execute_script("window.open('{}');".format(link))
        
        res = []
        for i in range(1, len(links)+1):
            driver.switch_to.window(driver.window_handles[i])
            curr = driver.current_url
            print(" >> getting link for {}".format(curr))
            if (len(curr) < 40): print(" >> failed to get link")
            res.append(driver.current_url)
        return res
    
    def close(self):
        self.driver.close()

if __name__ == "__main__":
    f = open("links/raw.txt", "r")
    temp = f.read().splitlines()
    f.close()

    print("=" * 100)    
    print(" + starting process")
    linker = Linker()
    #new_links = linker.get_links(temp)
    new_links = linker.fast_links(temp)
    linker.close()
    print("=" * 100)
    
    f = open("links/results.txt", "w")
    for link in new_links: f.write(link + "\n")
    f.close()