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

# Function to generate an INSERT query statement: keep querying on your endpoint to generate more statements
# Replicate these statements according to what we need to load into the database
def generate_insert_query():
#     For clickhouse
#     CREATE TABLE htap_table
#     (
#         user_name String,
#         user_age UInt8
#     ) engine=MergeTree ORDER BY user_age
#
#     For postgres
#     CREATE TABLE htap_table
#     (
#         id serial primary key,
#         user_name VARCHAR(20),
#         user_age INT,
#         user_email VARCHAR(50),
#         user_city VARCHAR(20),
#         user_country VARCHAR(20),
#         registration_date VARCHAR(20),
#         is_active BOOLEAN,
#         user_interest VARCHAR(20),
#         visit_count INT,
#         user_category VARCHAR(20)
#     )

    query = """
    INSERT INTO htap_table
    (
        user_name,
        user_age,
        user_email,
        user_city,
        user_country,
        registration_date,
        is_active,
        user_interest,
        visit_count,
        user_category) VALUES """
    values = []
    
    for _ in range(10):
        user_name = f"User_{random.randint(1, 10000)}"
        user_age = random.randint(18, 60)
        user_email = f"{user_name.lower().replace(' ', '_')}@example.com"
        user_city = f"City_{random.randint(1, 10)}"
        user_country = random.choice(["Country1", "Country2", "Country3"])
        registration_date = f"2022-01-{random.randint(1, 31)}"
        is_active = random.choice(["TRUE", "FALSE"])
        user_interest = f"Interest_{random.randint(1, 5)}"
        visit_count = random.randint(1, 100)
        user_category = random.choice(["Category1", "Category2", "Category3"])

        values.append(f"('{user_name}', {user_age}, '{user_email}', '{user_city}', '{user_country}', '{registration_date}', {is_active}, '{user_interest}', {visit_count}, '{user_category}')")

    return query + ','.join(values)


url = 'http://127.0.0.1:3333'

while True:
    time.sleep(0.1)
    q = generate_insert_query()
    x = requests.post(url+'/write', data = {'query': q})
