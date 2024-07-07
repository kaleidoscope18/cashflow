# Cashflow

Ledger application that allows you to know if your account will be overdraft someday.  
All features for this app are described in files under `bdd/features/`.

## Structure

The app will follow the DDD project architecture.  
Initial bootstrap is done in `main.go`, which is the entry file.  
It will instantiate the **repository** instances from `repository` package.  
Repositories are injected into the **services** of `domain` package.  
The services are used by **resolvers** (GraphQL) in `api` package.  
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

## Developing

Prerequisite: Docker

Open a terminal window and run

```sh
docker compose -f ./dev/docker-compose.yml --env-file=.env build            # for initial build, or if you change environment variables, docker or docker compose configs
docker compose -f ./dev/docker-compose.yml --env-file=.env up               # run the stack in containers, with debug
docker compose -f ./dev/docker-compose.yml --env-file=.env down             # stop / remove containers
```

You can enable/disable watch mode by using the "w" key while you develop.
Files that are not in .dockerignore will be watched.
Endpoints URLs will be provided in the logs.

### GraphQL

When changes are made in `.graphqls` files, you will need to regenerate code with graphqlgen by using this command:

```sh
cd api && rm -rf /graph/generated && go run github.com/99designs/gqlgen generate && cd ..
```

- All files under `api/graph/generated` will be regenerated
- If the regeneration throws an error, check your schemas first and backup the resolver files and erase their contents for a clean slate and try again.

### Debugging

When the containers are running in docker, launch `Go Containerized Debug` VSCode configuration.

### BDD

BDD testing is in `bdd` package and gherkin natural language is used with `features/*.feature` files.  
[Godog](https://github.com/cucumber/godog/) library is used to run these tests.

```sh
docker compose -f ./dev/docker-compose.yml -f ./dev/docker-compose.bdd.yml --env-file=.env up
```

#### Running BDD locally

You can run bdd tests locally when the app is (already) running on 8080.
Click on run / debug above the function named `TestBDD` if you have VSCode with the Go plugin.  

### Local commands

```sh
go test ./... -run 'Test[^BDD]'                 # run unit tests
go test ./... -run 'TestBDD'                    # run bdd tests
go get                                          # install dependencies
go mod tidy                                     # cleanup dependencies
go build                                        # build the binary
./cashflow                                      # run the binary
```

## Database

### In memory

For testing ONLY purposes, there's an in-memory strategy.

### MySQL

#### Run MySQL locally

1. `brew install mysql`
2. `brew services start mysql`
3. setup your database based on `dev/.seed.sql`, easy way to connect is via `mysql -u root`
4. edit your .env file to connect your app to your local db
5. run the app

## Infrastructure

The infra is on AWS written with Pulumi in Go.

### Deployment

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

Then you can connect to the mysql database with the pre-established credentials, i.e.:

```sh
mysql -h cashflow-db9b3ad36.c9gia6eo0ryf.us-east-1.rds.amazonaws.com -u admin -p
```

type `exit` to exit all instances.

## TODO

- e2e should maybe become bdd tests, e2e should only apply when there's also a front-end (playwright)
- add request context with tracing through the whole layers
- setup robust logging
