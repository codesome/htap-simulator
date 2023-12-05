# HTAP Simulator

<img src="assets/img/htap-simulator-architecture.png" alt="Architecture" width="400"/>


1) Start Postgres server and Clickhouse server on their default ports
2) Create a database named `htap` in both of them
3) Create the following tables

    In Postgres
    ```
    CREATE TABLE htap_table
    (
        id serial primary key,
        user_name VARCHAR(20),
        user_age INT,
        user_email VARCHAR(50),
        user_city VARCHAR(20),
        user_country VARCHAR(20),
        registration_date VARCHAR(20),
        is_active BOOLEAN,
        user_interest VARCHAR(20),
        visit_count INT,
        user_category VARCHAR(20)
    )
    ```

   In Clickhouse
    ```
    CREATE TABLE htap_table
    (
        user_name String,
        user_age UInt8
    )
    engine=MergeTree
    ORDER BY user_age
    ```

4) Start the htap-simulator

    ```bash
    $ go run htap-brain/*
    ```

5) Start the write and read load in parallel

    ```bash
    $ python3 load-generator/write.py
    $ python3 load-generator/read.py
    ```