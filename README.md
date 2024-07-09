# Cashflow

Ledger application that allows you to know if your account will be overdraft someday.  
All features for this app are described in files under `bdd/features/`.

## Structure

Packages:
| Name        | Description                                 |
| ---         | ---                                         |
| api         | graphql api                                 |
| bdd         | bdd tests in gherkin                        |
| dev         | dev tools (docker compose)                  |
| domain      | business logic                              |
| models      | contracts                                   |
| pulumi      | infra code                                  |
| repository  | database layer                              |
| utils       | reusable functions independent from domain  |

## Developing

Prerequisite: Docker

Open a terminal window and run

```sh
docker compose -f ./dev/docker-compose.yml --env-file=.env build    # for initial build, or if you change env variables or configs
docker compose -f ./dev/docker-compose.yml --env-file=.env up       # run the stack in containers, with debug
docker compose -f ./dev/docker-compose.yml --env-file=.env down     # stop / remove containers
```

You can enable/disable watch mode by using the "w" key while you develop.
Files that are not in .dockerignore will be watched.
Endpoints URLs will be provided in the logs.

### Development flow

I develop this app in a TDD way (red-green-refactor) + AI, here's how:  

- I always keep my dev containers running (see [Developing](#developing) above).  
- I first write the scenarios(features) I need to code in a `*.feature` file.  
- Then I run the BDD tests via to print the step defs I need to implement (see [BDD](#bdd)).  
- I implement them to a `*.stepdef.go` file.  
    The test will be failing first (red).
- I implement the code to make the tests pass (green).  
- I improve and refactor the code using Codium AI.
- I add unit tests automatically using Codium AI and Perplexity AI.
- I make sure the tests are not flaky by running the whole test suite many times.

### GraphQL

When changes are made in `.graphqls` files, you will need to regenerate code with graphqlgen by using this command:

```sh
cd api && rm -rf /graph/generated && go run github.com/99designs/gqlgen generate && cd ..
```

- All files under `api/graph/generated` will be regenerated
- If the regeneration throws an error, check your schemas first and backup the resolver files and erase their contents for a clean slate and try again.

### Debugging

When the containers are running in docker, launch `Go Containerized Debug` VSCode configuration.  

Fyi, here my installed VSCode Extensions:  

- Go by Go Team at Google  
- markdownlint  
- GraphQL: Syntax Highlighting  
- Prettier - Code formatter  
- Snippets and Syntax Highlight for Gherkin (Cucumber) by Euclidity  

### BDD

BDD testing is in `bdd` package and gherkin natural language is used with `features/*.feature` files.  
[Godog](https://github.com/cucumber/godog/) library is used to run these tests.

```sh
docker compose -f ./dev/docker-compose.yml -f ./dev/docker-compose.bdd.yml --env-file=.env up
```

#### Running BDD locally

You can run bdd tests locally when the app is (already) running on 8080.  
Click on run / debug above the function named `TestBDD` if you have VSCode with the Go plugin (TO REVIEW, SOMETIMES FALSE GREEN???).  

### Local commands

```sh
go test ./... -run 'Test[^BDD]'                 # run unit tests
go test ./... -run 'TestBDD'                    # run bdd tests
go get                                          # install dependencies
go mod tidy                                     # cleanup dependencies
go build                                        # build the binary
./cashflow                                      # run the binary
```

## Infrastructure

The infra is on AWS written with Pulumi in Go.

### Deployment

(login to AWS first)

```sh
cd pulumi/
pulumi stack select <dev | prod>
pulumi up
```

## Architecture

The app will follow the DDD project architecture.  
Initial bootstrap is done in `main.go`, which is the entry file.  
It will instantiate the **repository** instances from `repository` package.  
Repositories are injected into the **services** of `domain` package.  
The services are used by **resolvers** (GraphQL) in `api` package.  
All the domain logic is in `domain` package alongside services that use them.  
All domain models are in `models` package.  
For testing ONLY purposes, there's an in-memory strategy, but the app connects to a MySQL database.  
