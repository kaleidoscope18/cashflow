# Cashflow

Ledger application that allows you to know if your account will be overdraft someday.  
All features for this app are described in files under `bdd/features/`.
Developing with a prod database since this is a hobby project.

## Technologies

- Go
- GraphQL
- PostgreSQL
- Supabase

## Structure

Packages:

| Name        | Description                                 |
| ---         | ---                                         |
| api         | graphql api                                 |
| bdd         | bdd tests in gherkin                        |
| dev         | dev tools                                   |
| domain      | business logic                              |
| models      | contracts                                   |
| repository  | database layer                              |
| utils       | reusable functions independent from domain  |

## Commands

```sh
go get                                          # install dependencies
go mod tidy                                     # cleanup dependencies
go fmt                                          # format code
go build                                        # build the binary
./cashflow                                      # run the binary

# run tests
go test ./... -run 'Test[^BDD]'                 # run unit tests
go test ./... -run 'TestBDD'                    # run all bdd tests
go test ./... -run 'TestBDD' -- --tags=only,wip # run bdd tests with only and wip tags (i.e. @only above the scenario)
```

## Debug

- use `Debug Go` VSCode configuration to debug the app.
- click on `Debug test` above `TestBDD` in `bdd_test.go` file to debug the bdd tests.

## Development flow

I develop this app in a TDD way (red-green-refactor), here's how:  

- I first write the scenarios(features) I need to code in a `*.feature` file.  
- Then I run the BDD tests via to print the step defs I need to implement (see [BDD](#bdd)).  
- I implement them to a `*.stepdef.go` file.  
    The test will be failing first (red).
- I implement the code to make the tests pass (green).  
- I improve and refactor the code.
- I add unit tests with AI to speed me up.
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

### Suggested vscode extensions

- Go by Go Team at Google  
- markdownlint  
- GraphQL: Syntax Highlighting  
- Prettier - Code formatter  
- Snippets and Syntax Highlight for Gherkin (Cucumber) by Euclidity  

### BDD

BDD testing is in `bdd` package and gherkin natural language is used with `features/*.feature` files.  
[Godog](https://github.com/cucumber/godog/) library is used to run these tests.

## Architecture

The app will follow the DDD project architecture.  
Initial bootstrap is done in `main.go`, which is the entry file.  
It will instantiate the **repository** instances from `repository` package.  
Repositories are injected into the **services** of `domain` package.  
The services are used by **resolvers** (GraphQL) in `api` package.  
All the domain logic is in `domain` package alongside services that use them.  
All domain models are in `models` package.
