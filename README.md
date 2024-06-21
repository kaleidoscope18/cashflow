# Cashflow

Ledger application that allows you to know if your account will be overdraft someday. 

- Transaction and balance tracking for different accounts.
- Server runs GraphQL API.

## TODO

- e2e should maybe become bdd tests
- e2e should only apply when there's also a front-end (playwright)
- add request context with tracing through the whole layers
- setup robust logging

## Structure

Initial bootstrap is done in `main.go`, which is the entry file.
Depending on the storage strategy, it will instantiate the correct repository instances from `repository` package, and then inject them into the newly created services from `domain` package.
The services are used by api handlers (REST) or resolvers (GraphQL) in `api` package.
All the domain logic is in `domain` package alongside services that use them.
All domain models are in `models` package.

Packages:
| Name        | Description                       |
| ---         | ---                               |
| domain      | all business logic                |
| models      | all domain models                 |
| api         | graphql api                       |
| e2e         | e2e tests                         |
| repository  | infrastructure                    |
| utils       | functions independent from domain |
| bdd         | bdd tests - gherkin language      |

## Getting started

1. `go get`
2. `go mod tidy`
3. `go test ./...`
3. `go build`
4. `./cashflow`

Endpoints URLs will be provided by the running process logs.

## Debugging

Follow this : https://github.com/golang/vscode-go/wiki/debugging

## Features

This project's whole functionnality set is documented in Gherkin natural language. You can find them in the `*.feature` files under `bdd/`.

## Graphql (graphqlgen)

When changes are made in `.graphqls` files, you will need to regenerate code by using this command:

```sh
cd api && rm -rf /graph/generated && go run github.com/99designs/gqlgen generate && cd ..
```

- All files under `api/graph/generated` will be regenerated
- If the regeneration throws an error, check your schemas first and backup the resolver files and erase their contents for a clean slate and try again.

## BDD

BDD testing is in `bdd` package and gherkin natural language is used with `features/*.feature` files.
[Godog](https://github.com/cucumber/godog/) library is used to run these tests.

These tests are included in test run with `go test ./...` command from root.

## Database

### In memory

For testing purposes, there's an in-memory strategy and for production, this app is set to use a mysql db.

### MySQL

1. `brew install mysql`
2. `brew services start mysql`
3. `mysql -u root` (to quit, type `QUIT`)
4. edit database connection string `username:password@tcp(host:port)/database_name`

```sh
mysql> SHOW DATABASES;
mysql> CREATE DATABASE cashflow;                        # to remove: DROP DATABASE cashflow;
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

#### Connect locally

```sh
mysql -u root -pnew_password cashflow
```

## Infra

The infra is on AWS written in pulumi.

### Deployment

See pulumi's documentation for more details.

(login to AWS first)

```sh
cd pulumi/
pulumi stack select <environment>
pulumi up
```

### SSH

To connect to the public ec2 instance (bastion host)
```sh
ssh -i ~/.ssh/id_rsa ec2-user@<ec2 instance public ip>
```

Install mysql client with:
```sh
sudo su -
dnf -y localinstall https://dev.mysql.com/get/mysql80-community-release-el9-4.noarch.rpm
dnf -y install mysql mysql-community-client
exit
```

Optional: verify rds endpoint
```sh
sudo yum install -y nc
nc -zv <rds endpoint> 3306
```

Then you can connect to the mysql database with the pre-established credentials
```sh
mysql -h cashflow-db9b3ad36.c9gia6eo0ryf.us-east-1.rds.amazonaws.com -u admin -p
```

type `exit` to exit all instances.

### Docker

```sh
docker build -t kaleidoscope18/cashflow:1 ./ # docker build -t <USER>/<CONTAINER>:<VERSION> ./
docker run -p 8080:8080 kaleidoscope18/cashflow:1
```

Note: to cleanup old docker images and containers do `docker rm $(docker ps -a -q) && docker rmi $(docker images -a -q)`