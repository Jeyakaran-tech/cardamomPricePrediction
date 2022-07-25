# cardamomPricePrediction

This is an idea about web scraping the data from the site http://www.indianspices.com/indianspices/marketing/price/domestic/daily-price-large.html and use that data to predict the Cardamom Prices in the next 30 days using GCP's tools. The above mentioned site is about our local farmers cardamom prices. I really think predicting the prices would help them to decide on their work like when to harvest, peak time, and coming up with labour strategies according to the prices. This will definitely increase their throughput and improve their lives. 

This will also gives me some idea about the overall usage of GCP.

## Current Stage of the project

Currently web service is deployed in Cloud Run and once the service is called, the json(Data scraped) can be pushed automatically to GCS(Google Cloud Storage). 

## Yet to be done

Pipeline needs to be constructed to push the data to BigQuery and with the help of MachineLearning, do the prediction and push the data again into Data Lab. 