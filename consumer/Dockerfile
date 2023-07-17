FROM golang:1.20-alpine AS build-stage

RUN apk --no-cache add ca-certificates
WORKDIR /btcApplication
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /btcApplication/consumer .
RUN chmod +x /btcApplication/consumer

FROM scratch

COPY --from=build-stage /btcApplication/consumer /btcApplication/consumer

EXPOSE 8080
ENTRYPOINT ["/btcApplication/consumer"]