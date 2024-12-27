CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    description VARCHAR(200) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    date DATE NOT NULL,
    recurrency VARCHAR(500) NOT NULL
);

CREATE TABLE balances (
    amount DECIMAL(10,2) NOT NULL,
    date DATE NOT NULL,
    PRIMARY KEY (date)
);