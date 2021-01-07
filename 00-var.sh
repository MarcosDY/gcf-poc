#!/bin/bash

# Name for GCloud funtion to deploy
FUNCTION_NAME="svid-function"
# Project where to deploy function
PROJECT_ID="YOUR_PROJECT_ID"
# GCloud region
REGION="us-central1"
# Secret ID 
SECRET_ID="function-svid"
# Secret name where to consume X509-SVID it is in format "projects/*/secrets/*/version/latest"
SECRET_NAME="projects/${PROJECT_ID}/secrets/${SECRET_ID}/versions/latest"
# Service accout used for function. It must have permissions to 'secretmanager.versions.access'
SERVICE_ACCOUNT="YOUR_SERVICE_ACCOUNT"
