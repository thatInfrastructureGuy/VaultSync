go build ../vaultsync.go

# Required Environment Variables
export  AZURE_TENANT_ID=ad4b7142-cac0-4b8b-9573-c02311c68a26
export 	AZURE_CLIENT_ID=d70a6390-35a9-443e-9e8e-805cd93aa68b
export 	AZURE_CLIENT_SECRET=a40c02e7-2dab-47b3-a470-65d6451d6f13
export  KVAULT=kvtestin

# Use metadata.namespace from downward API when running in cluster
export 	SECRETS_NAMESPACE=kv-namespace

# Required Locally. Not required when running in cluster
export  KUBECONFIG=/home/play/.kube/config

# Additional Options
# Azure KeyVault does not support Underscores. To use underscores; we use hyphers ("dashes") in KeyVault and set the flag.
# Make sure to use only hyphens or only underscores in environment variables. Using both is NOT supported.
export 	CONVERT_HYPHENS_TO_UNDERSCORES=true

# Skip if you want SECRET_NAME = KVAULT
export  SECRET_NAME=kvtestin

# Debugs Keyvault requests/responses
export  DEBUG=false

./vaultsync
