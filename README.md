# Cashflow

## Structure



## Getting started

1. `go get`
2. `go mod tidy`
3. `go test ./...`
3. `go build`
4. `./cashflow`

## Debugging

Follow this : https://github.com/golang/vscode-go/wiki/debugging

## Features

- Add balance on a specified date
- Add transaction on a specified date
- List transactions

## Graphql (graphqlgen)

- All files under `api/graph/generated` will be regenerated
- You might have to backup the contents of the resolvers. Sometimes there are breaking changes.
- to regenerate do:

```sh
cd api && rm -rf /graph/generated && go run github.com/99designs/gqlgen generate
```

## Database

### MySQL

1. `brew install mysql`
2. `brew services start mysql`
3. `mysql -u root` (to quit, type `QUIT`)
4. edit database connection string `username:password@tcp(host:port)/database_name`

```sh
mysql> SHOW DATABASES;
mysql> CREATE DATABASE cashflow; # DROP DATABASE cashflow;
mysql> USE cashflow;
mysql> CREATE TABLE transactions
(
  id              INT unsigned NOT NULL AUTO_INCREMENT, # Unique ID for the record
  description     VARCHAR(200) NOT NULL,                # Transaction description
  amount          decimal(10,2) NOT NULL,               # Transaction amount
  date            DATE NOT NULL,                        # Date of the transaction
  PRIMARY KEY     (id)                                  # Make the id the primary key
);
mysql> SHOW TABLES;
mysql> DESCRIBE transactions;
mysql> CREATE TABLE balances
(
  amount          decimal(10,2) NOT NULL,               # Balance amount
  date            DATE NOT NULL,                        # Balance date
  PRIMARY KEY     (date)                                # Make the id the primary key
);
```