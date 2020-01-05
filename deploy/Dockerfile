FROM golang:1.13-alpine as build-env
MAINTAINER Ashish <thatInfrastructureGuy@gmail.com>
RUN adduser -u 1000 -D -H scratchuser

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -trimpath -tags netgo -ldflags '-s -w -extldflags "-static"' \
    -o /app/publish/vaultsync \
    /app/cmd/standalone/main.go

# Minimal Image
FROM scratch
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /etc/passwd /etc/passwd
USER scratchuser
WORKDIR /app
COPY --from=build-env --chown=1000:1000 /app/publish /app
ENTRYPOINT ["/app/vaultsync"]
