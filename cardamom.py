import requests
from bs4 import BeautifulSoup
import csv


URL = "http://www.indianspices.com/indianspices/marketing/price/domestic/daily-price-large.html"
page = requests.get(URL)

soup = BeautifulSoup(page.content, "html.parser")

numberOfPagesString = soup.find("div", class_="text-warning").text.strip().split(" ")
totalNumberOfPages = int(numberOfPagesString[3])

for x in range(1, 3):
    page = requests.get(URL+f"?page={x}")
    soup = BeautifulSoup(page.content, "html.parser")

    for rows in soup.find_all('tr')[5:]:
        for cols in rows.find_all('td'):
            print(cols.text.strip(), ",")
        print("\n")
       