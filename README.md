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