# Cashflow

Ledger application that allows you to know if your account will be overdraft someday. 

- Transaction and balance tracking for different accounts.
- Server runs GraphQL API.

## TODO

- e2e should maybe become bdd tests, e2e should only apply when there's also a front-end (playwright)
- add request context with tracing through the whole layers
- setup robust logging

## Structure

The app will follow the DDD project architecture.  
Initial bootstrap is done in `main.go`, which is the entry file.  
It will instantiate the **repository** instances from `repository` package.  
Repositories are injected into the **services** of `domain` package.  
The services are used by ~~api handlers (REST) or~~ **resolvers** (GraphQL) in `api` package.  
All the domain logic is in `domain` package alongside services that use them.  
All domain models are in `models` package.  

Packages:
| Name        | Description                                 |
| ---         | ---                                         |
| api         | graphql api                                 |
| bdd         | bdd tests in gherkin                        |
| dev         | dev tools                                   |
| domain      | business logic                              |
| models      | contracts                                   |
| pulumi      | infra code                                  |
| repository  | database layer                              |
| utils       | reusable functions independent from domain  |

## Getting started

Prerequisite: Docker

Open a terminal window and run
```sh
docker compose -f ./dev/docker-compose.yml --env-file=.env up
```

and then you can enable/disable watch mode by using the "w" key while you develop.

Alternative: open two terminal windows and run these commands, one each
```sh
docker compose -f ./dev/docker-compose.yml --env-file=.env watch
docker compose -f ./dev/docker-compose.yml logs -f
```

Endpoints URLs will be provided in the logs.

To shut down containers:

```sh
docker compose -f ./dev/docker-compose.yml --env-file=.env down
```

You might have to rebuild if you change docker compose's configs

```sh
docker compose -f ./dev/docker-compose.yml --env-file=.env build
```

### Developing locally

1. Setup Mysql = go to [Run MySQL locally](#run-mysql-locally)
2. `go get`
3. `go mod tidy`
4. `go test ./...`
5. `go build`
6. `./cashflow`

## Debugging

Follow this : https://github.com/golang/vscode-go/wiki/debugging

## Features

This project's whole functionality set is documented in Gherkin natural language. You can find them in the `*.feature` files under `bdd/`.

### BDD

BDD testing is in `bdd` package and gherkin natural language is used with `features/*.feature` files.  
[Godog](https://github.com/cucumber/godog/) library is used to run these tests.

These tests are included in test run with `go test ./...` command from root.

You can run bdd only:
```sh
cd bdd/
go test -test.v -test.run ^TestFeatures$
```

You can run one scenario at a time  
(tip: previous command will list all scenarios)

```sh
go test -test.v -test.run ^TestFeatures/Adding_a_recurring_transaction$
```

## Graphql (graphqlgen)

When changes are made in `.graphqls` files, you will need to regenerate code by using this command:

```sh
cd api && rm -rf /graph/generated && go run github.com/99designs/gqlgen generate && cd ..
```

- All files under `api/graph/generated` will be regenerated
- If the regeneration throws an error, check your schemas first and backup the resolver files and erase their contents for a clean slate and try again.

## Database

### In memory

For testing purposes, there's an in-memory strategy.

### MySQL

#### Run MySQL locally

1. `brew install mysql`
2. `brew services start mysql`
3. edit your .env file and make sure you got the correct database connection string `username:password@tcp(host:port)/database_name`

You can perform all sorts of operations on the database with `mysql -u root` (to quit, type `QUIT`) 

#### Connect locally

```sh
mysql -u root -p<root-password> cashflow
```

## Infra

The infra is on AWS written with Pulumi in Go.

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
