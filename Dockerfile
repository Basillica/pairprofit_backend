##
## Build
##
FROM golang:1.17.2-buster AS build

WORKDIR /app
COPY . .

COPY .git ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY helpers/*.go ./helpers/
# handlers
COPY handlers/auth/*.go ./handlers/auth/

COPY middleware/*.go ./middleware/
COPY types/*.go ./types/
# COPY .env ./
RUN export GIT_COMMIT=$(git rev-list -1 HEAD) && \
  go build -ldflags "-X main.GitVersion=$GIT_COMMIT" -o /pairprofit_backend

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /pairprofit_backend /pairprofit_backend

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/pairprofit_backend"]