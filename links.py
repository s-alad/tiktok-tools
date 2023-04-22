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
            new_links.append(self.get_link(link))
        return new_links
    
    def close(self):
        self.driver.close()

if __name__ == "__main__":
    f = open("links/test.txt", "r")
    temp = f.read().splitlines()
    f.close()

    linker = Linker()
    new_links = linker.get_links(temp)
    linker.close()
    
    f = open("links/results.txt", "w")
    for link in new_links: f.write(link + "\n")
    f.close()