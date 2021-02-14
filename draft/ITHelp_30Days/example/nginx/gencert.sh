#!/bin/bash

openssl genrsa -out nginx-uat.key 2048

openssl req -new -key nginx-uat.key -out nginx-uat.csr  \
  -subj "/C=US/ST=CA/L=Mountain View/O=OS3/OU=Eng/CN=nginx-uat.apps-crc.testing"

openssl x509 -req -days 366 -in nginx-uat.csr  \
      -signkey nginx-uat.key -out nginx-uat.crt
