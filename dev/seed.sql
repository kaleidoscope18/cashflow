USE cashflow;

CREATE TABLE transactions
(
  id              INT unsigned NOT NULL AUTO_INCREMENT, -- Unique ID for the record
  description     VARCHAR(200) NOT NULL,                -- Transaction description
  amount          decimal(10,2) NOT NULL,               -- Transaction amount
  date            DATE NOT NULL,                        -- Date of the transaction
  recurrency      VARCHAR(500) NOT NULL,                -- Recurrency expression
  PRIMARY KEY     (id)                                  -- Make the id the primary key
);

CREATE TABLE balances
(
  amount          decimal(10,2) NOT NULL,               -- Balance amount
  date            DATE NOT NULL,                        -- Balance date
  PRIMARY KEY     (date)                                -- Make the id the primary key
);