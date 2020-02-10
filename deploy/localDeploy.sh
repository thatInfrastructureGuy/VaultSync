###############################################################################

# Copyright 2020 Kulkarni, Ashish <thatInfrastructureGuy@gmail.com>
# Author: Ashish Kulkarni
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

###############################################################################
# Build the binary
VERSION=$(git describe --tags --always --dirty)
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
CGO_ENABLED=0 go build -a -trimpath -tags netgo -ldflags "-s -w -extldflags '-static' \
    -X main.CodeVersion=$VERSION -X 'main.BuildTime=$BUILD_DATE' -X 'main.GoVersion=`go version`'" \
    -o vaultsync cmd/standalone/main.go

exit 1

###############################################################################
# Required Environment Variables

# Azure Credentials
export AZURE_TENANT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export AZURE_CLIENT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export AZURE_CLIENT_SECRET=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

# AWS Credentials
export AWS_ACCESS_KEY_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export AWS_SECRET_ACCESS_KEY=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export AWS_DEFAULT_REGION=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export AWS_REGION=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx


# KeyVault / Secrets Manager Name
export VAULT_NAME=mykeyvault

###############################################################################
# Optional Environment Variables

export CONSUMER=kubernetes
# Required Locally. Not required when running in cluster
export KUBECONFIG=~/.kube/config
# Skip if you want SECRET_NAME = VAULT_NAME
export SECRET_NAME=mysecret
export SECRET_NAMESPACE=mynamespace

###############################################################################
# Refresh rate syncs periodically.
# For one-time run such as running as init container set it to 0.
# Defaults to 60 // in Seconds
export REFRESH_RATE=0

###############################################################################
# Additional Options
# Azure KeyVault does not support Underscores. To use underscores; we use hyphers ("dashes") in KeyVault and set the flag.
# Make sure to use only hyphens or only underscores in environment variables. Using both is NOT supported.
export CONVERT_HYPHENS_TO_UNDERSCORES=true

# Debugs Keyvault requests/responses
export DEBUG=false

###############################################################################
./vaultsync
