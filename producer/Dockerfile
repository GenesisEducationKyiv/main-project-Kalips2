FROM golang:1.20-alpine AS build-stage

RUN apk --no-cache add ca-certificates
WORKDIR /btcApplication
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /btcApplication/producer .
RUN chmod +x /btcApplication/producer

FROM scratch

COPY --from=build-stage /btcApplication/producer /btcApplication/producer
COPY --from=build-stage /btcApplication/subscriptions.csv ./
COPY --from=build-stage /btcApplication/config.json/ ./
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080
ENTRYPOINT ["/btcApplication/producer"]