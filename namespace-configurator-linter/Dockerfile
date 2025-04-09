FROM golang:alpine3.17 AS build
ENV  CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY config ./config
RUN GOOS=linux GOARCH=amd64 go build -o ./namespace-configurator-linter

FROM  golang:alpine3.17
WORKDIR /app
COPY --from=build /app/namespace-configurator-linter ./
COPY --from=build /app/config ./config
ENTRYPOINT ["/app/namespace-configurator-linter"]
