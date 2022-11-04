#!/bin/sh

# Don't try any existing creds
export AWS_ACCESS_KEY_ID=fake
export AWS_SECRET_ACCESS_KEY=fake
export AWS_SESSION_TOKEN=fake

# Don't attempt EC2 IMD endpoint
export AWS_EC2_METADATA_DISABLED=true
export AWS_CONTAINER_CREDENTIALS_FULL_URI=http://999.999.999.999:1/
export AWS_CONTAINER_AUTHORIZATION_TOKEN="Basic YnJva2VuCg=="

# Don't load from any CLI files
export AWS_CONFIG_FILE=fake
export AWS_SHARED_CREDENTIALS_FILE=fake
# export AWS_DEFAULT_REGION=us-east-1

exec "$@"