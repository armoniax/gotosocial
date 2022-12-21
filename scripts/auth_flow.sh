#!/bin/sh

set -eux

SERVER_URL="http://localhost:8080"
#REDIRECT_URI="${SERVER_URL}"
#CLIENT_NAME="Test Application Name"
#REGISTRATION_REASON="Testing whether or not this dang diggity thing works!"
#REGISTRATION_USERNAME="${1}"
##REGISTRATION_EMAIL="${2}"
#REGISTRATION_PASSWORD=${2}
#REGISTRATION_EMAIL=${REGISTRATION_PASSWORD}"@amax.com"
#REGISTRATION_AGREEMENT="true"
#REGISTRATION_LOCALE="en"
#
## Step 1: create the app to register the new account
#CREATE_APP_RESPONSE=$(curl --fail -s -X POST -F "client_name=${CLIENT_NAME}" -F "redirect_uris=${REDIRECT_URI}" "${SERVER_URL}/api/v1/apps")
#CLIENT_ID=$(echo "${CREATE_APP_RESPONSE}" | jq -r .client_id)
#CLIENT_SECRET=$(echo "${CREATE_APP_RESPONSE}" | jq -r .client_secret)
#echo "Obtained client_id: ${CLIENT_ID} and client_secret: ${CLIENT_SECRET}"
#
## Step 2: obtain a code for that app
#APP_CODE_RESPONSE=$(curl --fail -s -X POST -F "scope=read" -F "grant_type=client_credentials" -F "client_id=${CLIENT_ID}" -F "client_secret=${CLIENT_SECRET}" -F "redirect_uri=${REDIRECT_URI}" "${SERVER_URL}/oauth/token")
#APP_ACCESS_TOKEN=$(echo "${APP_CODE_RESPONSE}" | jq -r .access_token)
#echo "Obtained app access token: ${APP_ACCESS_TOKEN}"
#
## Step 3: use the code to register a new account
#ACCOUNT_REGISTER_RESPONSE=$(curl --fail -s -H "Authorization: Bearer ${APP_ACCESS_TOKEN}" -F "reason=${REGISTRATION_REASON}" -F "email=${REGISTRATION_EMAIL}" -F "username=${REGISTRATION_USERNAME}" -F "password=${REGISTRATION_PASSWORD}" -F "agreement=${REGISTRATION_AGREEMENT}" -F "locale=${REGISTRATION_LOCALE}" "${SERVER_URL}/api/v1/accounts")
#USER_ACCESS_TOKEN=$(echo "${ACCOUNT_REGISTER_RESPONSE}" | jq -r .access_token)
#echo "Obtained user access token: ${USER_ACCESS_TOKEN}"
#
## # Step 4: verify the returned access token
#VERIFY_RESPONSE=$(curl -s -H "Authorization: Bearer ${USER_ACCESS_TOKEN}" "${SERVER_URL}/api/v1/accounts/verify_credentials")
#echo "verify_credentials: ${VERIFY_RESPONSE}"

## Step 5: change email
#CHANGE_EMAIL=$(curl --fail -s -X POST -H "Authorization: Bearer ${USER_ACCESS_TOKEN}" -F new_email="-"  "${SERVER_URL}/api/v1/user/email_change")
#echo "email: ${CHANGE_EMAIL}"

# Step 6: Sign in
#USER_TOKEN=$(curl --fail -s -X POST  -F "username=${REGISTRATION_USERNAME}"  -F "pub_key=${REGISTRATION_PASSWORD}"  "${SERVER_URL}/auth/sign_in/unconfirmed_email")
#echo "signed in token: ${USER_TOKEN}"

# Step 7 amax create
#AMAX_CREATE=$(curl --fail -s -H "Authorization: Bearer NJRKMTG1MTCTMJK5NC0ZOTE4LTLIZTKTNZU0MDRKNDMWYJLM" -F "username=jackwang54" -F "pub_key=biancheng347C12345679Abcdefg54" -F "client_id=01HTYTW1H0TG8GMGFW93MXKR1Y" -F "redirect_uri=http://localhost:8080" -F "scopes=read"  -F "response_type=json" "${SERVER_URL}/api/v1/accounts/submit_amax_info")
#echo "Obtained user access token: ${AMAX_CREATE}"

# Step 8 login for username and pubkey
USER_TOKEN=$(curl --fail -s -X POST  -F "username=jackwang54"  -F "pub_key=biancheng347C12345679Abcdefg54"  "${SERVER_URL}/oauth/token/unconfirmed_email")
echo "signed in token: ${USER_TOKEN}"

