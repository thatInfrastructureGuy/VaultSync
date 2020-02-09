###############################################################################
# Build the binary
CGO_ENABLED=0 go build -a -trimpath -tags netgo -ldflags '-s -w -extldflags "-static"' \
    -v -o vaultsync cmd/standalone/main.go

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
