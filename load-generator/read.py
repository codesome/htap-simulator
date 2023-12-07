"""
Spits out SELECT statements to be used to load data into a database

Select selects a random number of columns from the table
    Examples:
        SELECT user_city, user_email, user_interest, user_category, user_name, user_age FROM your_table;
        SELECT user_city, user_country, user_name, registration_date FROM your_table;
        SELECT visit_count FROM your_table;

    With probability of 0.5, it also gives out a fixed OLAP query for avg(user_age).
"""

import http.server
import socketserver
import json
import random
import requests
import time

# Function to generate and return a SELECT statement based on a random number
def generate_select_statement():
    # Randomly have OLAP query
    # TODO: NOTE: Remove this True to generate OLTP queries as well.
    if True or random.randint(0,1) == 1:
        return "select AVG(user_age) from htap_table"

    random_number = random.randint(1, 10)

    # Can add another list depending on pur schema
    values_random = [
        "user_name",
        "user_age",
        "user_email",
        "user_city",
        "user_country",
        "registration_date",
        "is_active",
        "user_interest",
        "visit_count",
        "user_category"
    ]
    # Randomly select random_number of columns from random positions
    selected_columns = random.sample(values_random, random_number)

    query = "SELECT * FROM htap_table"
    query = query.replace("*", ", ".join(selected_columns))
    query = query + f" LIMIT {random.randint(500, 1000)}"

    return query

url = 'http://127.0.0.1:3333'

while True:
    time.sleep(1)
    q = generate_select_statement()
    x = requests.get(url+'/read', params = {'query': q})
