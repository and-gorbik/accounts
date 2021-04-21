-- primary keys
ALTER TABLE city ADD PRIMARY KEY (id);
ALTER TABLE country ADD PRIMARY KEY (id);
ALTER TABLE account ADD PRIMARY KEY (id);

-- unique constraints
ALTER TABLE person ADD CONSTRAINT unique_account_id UNIQUE (account_id);
ALTER TABLE person ADD CONSTRAINT unique_email UNIQUE (email);


-- foreign keys
ALTER TABLE interest ADD FOREIGN KEY (account_id) REFERENCES account (id);
ALTER TABLE likes ADD FOREIGN KEY (liker_id) REFERENCES account (id);
ALTER TABLE likes ADD FOREIGN KEY (likee_id) REFERENCES account (id);
ALTER TABLE person ADD FOREIGN KEY (account_id) REFERENCES account (id);
ALTER TABLE person ADD FOREIGN KEY (country_id) REFERENCES country (id);
ALTER TABLE person ADD FOREIGN KEY (city_id) REFERENCES city (id);

-- indexes