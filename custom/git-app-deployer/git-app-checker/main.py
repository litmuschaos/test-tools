import base64
from random import randint, choice
import json
import os
import requests 
import time
import uuid
from itertools import count

""" Variables to maintain customer ID and header """
cust_id = ""
base64_str = base64.encodebytes(('%s:%s' % ("user","password")).encode()).decode().strip()

""" URL of the api endpoints of application to send and recieve requests """
url = os.environ['URL']

def check_for_user_login():
    try:        
        login = requests.get(url  +  "/login", headers={"Authorization":"Basic %s" % base64_str}, timeout=1)
        with login as response:
            if response.status_code == 200:
                print("[Status]: ResponseCode:", response.status_code," User successfully logged in")
            else:
                print("[Error]: ResponseCode:", response.status_code," User failed to logged in")    
                
        cust_id = login.cookies["logged_in"]
        cookies = {'logged_id': cust_id}
        print("[Info]: Unique ID for user is:", cust_id)
        
        #front-end
        with requests.get(url  +  "/", timeout=1) as response:
            print("[Info]: Sending request to front-end")
            if response.status_code == 200:
                print("[Status]: ResponseCode:", response.status_code," FrontEnd is accessible")
            else:
                print("[Error]: ResponseCode:", response.status_code," Front-End is unavailable")

        #register 
        username = "test_user_" + str(uuid.uuid4())
        password = "test_password"

        with requests.post(url  +  "/register", json={"username": username, "password": password}) as response:
            print("[Info]: Adding new cutomer ")
            if response.status_code == 200:
                print("[Status]: ResponseCode:", response.status_code, " Customer added successfully with user Name:",  username)
            else:
                print("[Error]: ResponseCode:", response.status_code, " Failed to add Customer with user Name:",  username)
    
    except Exception as exc:
        print("[Error]: Unable to send requests to used-db, Server encountered a condition that is preventing request.")
        print('[Error]: {err}'.format(err=exc))

def check_for_catalogue_items():
    try:
        #catalogue
        catalogue = requests.get(url  +  "/catalogue", timeout=1)
        with catalogue as response:
            if response.status_code == 200:
                category_item = choice(catalogue.json())
                print("[Status]: ResponseCode:", response.status_code," Catalogue get request successfully send")
                print("[Info]: Catalogue Item:", category_item)
            else:
                print("[Error]: ResponseCode:", response.status_code," Failed to get catalogue items")
    except Exception as exc:
        print("[Error]: Unable to send requests to calalogue-db, Server encountered a condition that is preventing request.")
        print('[Error]: {err}'.format(err=exc))

def user_details():
    try:
        #cards
        with requests.post(url  +"/cards", json={"longNum": "5429804235432", "expires": "04/16", "ccv": "432", "userId": cust_id},headers={"Authorization":"Basic %s" % base64_str}, timeout=1) as response:
            print("[Info]: Adding card details for purchase")
            if response.status_code == 200:
                print("[Status]: ResponseCode:", response.status_code," Card details has been successfully added")
            else:
                print("[Error]: ResponseCode:", response.status_code," Failed to add Card details")
        
        #addresses
        print("[Info]: Adding Address for user")
        with requests.post(url  +  "/addresses", json={"street": "my road", "number": "3", "country": "UK", "city": "London"},headers={"Authorization":"Basic %s" % base64_str}, timeout=1) as response:
            if response.status_code == 200:
                print("[Status]: ResponseCode:", response.status_code," Address has been added successfully")
            else:
                print("[Error]: ResponseCode:", response.status_code," Failed to add Address")
        
        #cards
        cards = requests.get(url  +"/cards", timeout=1)
        with cards as response:
            if response.status_code == 200:
                cardDetails = cards.json()
                card_id = choice(cardDetails["_embedded"]["card"])
                print("[Status]: RespondeCode:", response.status_code," Card ID:",card_id)
            else:
                print("[Error]: ResponseCode:", response.status_code," Failed to get Cards")
        
        #addresses
        with requests.get(url  +  "/addresses", timeout=1) as response:
            print("[Info]: Getting Address")
            if response.status_code == 200:
                print("[Status]: ResponseCode:", response.status_code," Address has been retrieved successfully")
            else:
                print("[Error]: ResponseCode:", response.status_code," Failed to add Address")    
        
        #details
        with requests.get(url  +  "/category.html", timeout=1) as response:
            print("[Info]: Getting item details")
            if response.status_code == 200:
                print("[Status]: ResponseCode:", response.status_code," Item details retrieved successfully")
            else:
                print("[Error]: ResponseCode:", response.status_code," Failed to get Items")    

    except Exception as exc:
        print("[Error]: Unable to send requests to user-db, Server encountered a condition that is preventing request.")
        print('[Error]: {err}'.format(err=exc))
               

print("[Status]: API Checker has been started")
while True:
    check_for_user_login()
    check_for_catalogue_items()
    user_details()
    time.sleep(6)
