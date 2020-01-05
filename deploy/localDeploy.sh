# Build the binary
CGO_ENABLED=0 go build -a -trimpath -tags netgo -ldflags '-s -w -extldflags "-static"' \
    -v -o vaultsync cmd/standalone/main.go

# Required Environment Variables
export  AZURE_TENANT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export 	AZURE_CLIENT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export 	AZURE_CLIENT_SECRET=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export  KVAULT=mykeyvault

# Use metadata.namespace from downward API when running in cluster
export 	SECRETS_NAMESPACE=mynamespace

# Required Locally. Not required when running in cluster
export  KUBECONFIG=/home/johnsmith/.kube/config

# Additional Options
# Azure KeyVault does not support Underscores. To use underscores; we use hyphers ("dashes") in KeyVault and set the flag.
# Make sure to use only hyphens or only underscores in environment variables. Using both is NOT supported.
export 	CONVERT_HYPHENS_TO_UNDERSCORES=true

# Skip if you want SECRET_NAME = KVAULT
export  SECRET_NAME=mysecret

# Debugs Keyvault requests/responses
export  DEBUG=false

./vaultsync
