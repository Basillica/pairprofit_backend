PairProfit Backend
===
The  backend application for the PairProfit software.
## Description
Authentication is with simple jwt. Token expiration times are stored in Redis and subsequently considered during authentication to cater for validity of tokens when a users logs out.
## Environment file
Create a `.env` file based on the following template
```bash
AWS_REGION="eu-central-1"
COOKIE_DOMAIN="localhost"
APP_ENVIRONMENT="local"  #for local dev
COOKIE_SECURE_ENABLE=false
DECODING_SECRET="09WJPdT9BUBXhH4A"
ALLOWED_ORIGIN="http://localhost:3000"
COOKIE_HTTPONLY=true
HUBSPOT_ENABLE=false
```

## Running one of the two release modes with Make
```bash
make # for debug mode
make run-debug # for debug mode
make run-release # for release mode
```

## Running using Docker and environment variables
```bash
docker compose up --build
```

## VS Code
Add missing YAML tags in settings.json, otherwise you'll see validation errors in `triggers/template.yaml`.
```
"yaml.customTags": ["!Ref", "!GetAtt"]
```

## Postman testing
Import `postman_collection.json` and test the flow.

[GO celery](https://github.com/gocelery/gocelery/tree/master/example)

## Gin custom struct validators
Some important links regarding validators
1. [custom validators](https://github.com/gin-gonic/gin#custom-validators)
2. [Example of validators](https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/)
###### xSI7s_6WFfa1-TIZLyE5pg

## Protocol Buffers
```bash
    protoc -I=pb --go_out=. pb/*.proto
```

## Documentation
To be able to view the entire structure of the code as well the documentation, [godocs](https://pkg.go.dev/golang.org/x/tools/cmd/godoc) package is used. 
This is easily installed using 
```bash
    go install golang.org/x/tools/cmd/godoc@latest
```
This can then be served locally using any port of choice as follows
```bash
    godoc -http=:6060
```
Navigate then to the link `http://localhost:6060` in the browser to explore the docs as well as any added examples from the source code as well as from external packages.
