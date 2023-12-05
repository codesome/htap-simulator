"""
Spits out INSERT and SELECT statements to be used to load data into a database
Insert always inserts 10 columns data
    Example: INSERT INTO your_table (user_name, user_age, user_email, user_city, user_country, registration_date, is_active, user_interest,
    visit_count, user_category) VALUES ('User_3821', 50, 'user_3821@example.com', 'City_10', 'Country2', '2022-01-8', False, 'Interest_5', 6, 'Category2');

Select always selects a random number of columns from the table
    Examples:
        SELECT user_city, user_email, user_interest, user_category, user_name, user_age FROM your_table;
        SELECT user_city, user_country, user_name, registration_date FROM your_table;
        SELECT visit_count FROM your_table;
"""

import http.server
import socketserver
import json
import random
import requests
import time

# Function to generate and return a SELECT statement based on a random number
def generate_select_statement():
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
    query = query + f" LIMIT {random.randint(100, 1000)}"
    if random.randint(0,1) == 1:
        query = "select AVG(user_age) from htap_table"
    return query

url = 'http://127.0.0.1:3333'

while True:
    time.sleep(0.5)
    q = generate_select_statement()
    x = requests.get(url+'/read', params = {'query': q})
