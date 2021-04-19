-- primary keys
ALTER TABLE city ADD PRIMARY KEY (id);
ALTER TABLE country ADD PRIMARY KEY (id);
ALTER TABLE person ADD PRIMARY KEY (id);

-- foreign keys
ALTER TABLE interest ADD FOREIGN KEY (account_id) REFERENCES account (id);
ALTER TABLE like ADD FOREIGN KEY (liker_id) REFERENCES account (id);
ALTER TABLE like ADD FOREIGN KEY (likee_id) REFERENCES account (id);
ALTER TABLE person ADD FOREIGN KEY (account_id) REFERENCES account (id);
ALTER TABLE person ADD FOREIGN KEY (country_id) REFERENCES country (id);
ALTER TABLE person ADD FOREIGN KEY (city_id) REFERENCES city (id);

-- unique constraints
ALTER TABLE person ADD CONSTRAINT unique_account_id UNIQUE (account_id);
ALTER TABLE person ADD CONSTRAINT unique_email UNIQUE (email);
ALTER TABLE person ADD CONSTRAINT unique_country_city UNIQUE (country, city);

-- indexes