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

# Metadata params
ARG BUILD_DATE
ARG VERSION
ARG VCS_URL
ARG VCS_REF

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -trimpath -tags netgo -ldflags "-s -w -extldflags '-static' \
    -X main.CodeVersion=$VERSION -X 'main.BuildTime=$BUILD_DATE' -X 'main.GoVersion=`go version`'" \
    -o /app/publish/vaultsync \
    /app/cmd/standalone/main.go

# Minimal Image
FROM scratch
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /etc/passwd /etc/passwd
USER scratchuser
WORKDIR /app

# Metadata
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="VaultSync" \
      org.label-schema.description="Periodically syncs secrets from various Vaults to Kubernetes Secrets." \
      org.label-schema.url="https://thatInfrastructureGuy.com" \
      org.label-schema.vcs-url=$VCS_URL \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vendor="thatInfrastructureGuy" \
      org.label-schema.version=$VERSION \
      org.label-schema.schema-version="1.0" 

COPY --from=build-env --chown=1000:1000 /app/publish /app
ENTRYPOINT ["/app/vaultsync"]
