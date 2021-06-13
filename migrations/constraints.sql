-- primary keys
ALTER TABLE city ADD PRIMARY KEY (id);
ALTER TABLE country ADD PRIMARY KEY (id);
ALTER TABLE account ADD PRIMARY KEY (id);

-- unique constraints
ALTER TABLE account ADD CONSTRAINT unique_account_email UNIQUE (email);
ALTER TABLE city ADD CONSTRAINT unique_city_name UNIQUE (name);
ALTER TABLE country ADD CONSTRAINT unique_country_name UNIQUE (name);

-- foreign keys
ALTER TABLE interest ADD FOREIGN KEY (account_id) REFERENCES account (id);
ALTER TABLE likes ADD FOREIGN KEY (liker_id) REFERENCES account (id);
ALTER TABLE likes ADD FOREIGN KEY (likee_id) REFERENCES account (id);
ALTER TABLE account ADD FOREIGN KEY (country_id) REFERENCES country (id);
ALTER TABLE account ADD FOREIGN KEY (city_id) REFERENCES city (id);

-- indexes