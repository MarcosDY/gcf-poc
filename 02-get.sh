#!/bin/bash

. 00-var.sh

curl "https://${REGION}-${PROJECT_ID}.cloudfunctions.net/${FUNCTION_NAME}"
 
