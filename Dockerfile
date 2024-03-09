FROM golang:1.20.5-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o ./out/go-app .

FROM alpine:3.16.2
COPY --from=build-base /app/out/go-app /app/go-app
COPY --from=build-base /app/shared/vehicle/*.json /app/shared/vehicle/

CMD ["/app/go-app"]
