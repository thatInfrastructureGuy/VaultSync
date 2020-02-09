#  Copyright 2020 Kulkarni, Ashish <thatInfrastructureGuy@gmail.com>
#  Author: Ashish Kulkarni
# 
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
# 
#    http://www.apache.org/licenses/LICENSE-2.0
# 
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

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
