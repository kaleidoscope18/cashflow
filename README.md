# Cashflow

## Getting started

1. `go get`
2. `go mod tidy`
3. `go test ./...`
3. `go build`
4. `./cashflow`

## Features

- Add balance on a specified date
- Add transaction on a specified date
- List transactions

## Info

- All files under `api/graph/generated` will be regenerated
- You might have to backup the contents of the resolvers. Sometimes there are breaking changes.
- to regenerate do:

```sh
cd api && rm -rf /graph/generated && go run github.com/99designs/gqlgen generate
```