# cryptoshare

## Generate Private, Public Key Pair

generate key paris

- `openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048`

- `openssl rsa -in rsa_private.pem -pubout -out rsa_public.pem`

## Requirements

Infura key

- `get infura api key from infura.io`

IP2Location

- `download bin file from https://lite.ip2location.com/database-download and specifically IP-COUNTRY-REGION-CITY`

- `place in conf folder`
