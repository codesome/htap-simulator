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

# Function to generate an INSERT query statement: keep querying on your endpoint to generate more statements
# Replicate these statements according to what we need to load into the database
def generate_insert_query():
    query = "INSERT INTO your_table (user_name, user_age, user_email, user_city, user_country, registration_date, is_active, user_interest, visit_count, user_category) VALUES "
    values = []
    
    # for _ in range(10000):
    user_name = f"User_{random.randint(1, 10000)}"
    user_age = random.randint(18, 60)
    user_email = f"{user_name.lower().replace(' ', '_')}@example.com"
    user_city = f"City_{random.randint(1, 10)}"
    user_country = random.choice(["Country1", "Country2", "Country3"])
    registration_date = f"2022-01-{random.randint(1, 31)}"
    is_active = random.choice([True, False])
    user_interest = f"Interest_{random.randint(1, 5)}"
    visit_count = random.randint(1, 100)
    user_category = random.choice(["Category1", "Category2", "Category3"])

    values.append(f"('{user_name}', {user_age}, '{user_email}', '{user_city}', '{user_country}', '{registration_date}', {is_active}, '{user_interest}', {visit_count}, '{user_category}')")

    return query + ', '.join(values) + ';'

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

    query = "SELECT * FROM your_table;"
    query = query.replace("*", ", ".join(selected_columns)) 
    return query 


# Custom handler to respond to a different endpoints: can return as json if you want
class MyHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/generate_insert_query':
            self.send_response(200)
            self.send_header('Content-type', 'text/plain')
            self.end_headers()

            # Generate and send INSERT query statement
            insert_query = generate_insert_query()
            self.wfile.write(insert_query.encode())
        elif self.path == '/generate_select_query':
            self.send_response(200)
            self.send_header('Content-type', 'text/plain')
            self.end_headers()

            # Generate and send SELECT query statement
            select_query = generate_select_statement()
            self.wfile.write(select_query.encode())
        else:
            # Default behavior for other paths
            super().do_GET()

# Create the server that runs until broken from inside like Algo exams :D
PORT = 8000
with socketserver.TCPServer(("", PORT), MyHandler) as httpd:
    print(f"Serving on port {PORT}")
    httpd.serve_forever()
