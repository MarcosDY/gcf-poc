#!/bin/bash

. 00-var.sh

gcloud functions deploy $FUNCTION_NAME --source function --entry-point SvidGet --runtime go113 --trigger-http --allow-unauthenticated --service-account="${SERVICE_ACCOUNT}" --set-env-vars SECRET_NAME="${SECRET_NAME}"
 
