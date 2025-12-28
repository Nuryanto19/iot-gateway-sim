#!/bin/bash

# Target directory
GATEWAY_DIR="gateway-certs"
MOSQUITTO_DIR="infra/mosquitto/config/certs"

# Create new folder if not found
mkdir -p $GATEWAY_DIR
mkdir -p $MOSQUITTO_DIR

echo "Generating CA..."
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -sha256 -days 365 -out ca.crt \
    -subj "/C=ID/ST=Jakarta/L=Jakarta/O=MyOrg/CN=my-ca"

echo "Generating Server Cert (Mosquitto) with SAN..."
openssl genrsa -out server.key 2048
# Create temporary config file for SAN server
# DNS.2 : ensure match with container name in docker-compose.yaml
cat > server_ext.cnf <<EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no
[req_distinguished_name]
C = ID
ST = Jakarta
L = Jakarta
O = MyOrg
CN = localhost
[v3_req]
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
DNS.2 = iot-broker
IP.1 = 127.0.0.1
EOF

openssl req -new -key server.key -out server.csr -config server_ext.cnf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out server.crt -days 360 -sha256 -extensions v3_req -extfile server_ext.cnf

echo "Generating Client Cert (Go Gateway)..."
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr \
    -subj "/C=ID/ST=Jakarta/L=Jakarta/O=MyOrg/CN=iot-gateway-1"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out client.crt -days 350 -sha256

# Distributr file to target directory
echo "Distributing certificates..."

# For Go Gateway (Need CA, Client Cert, Client Key)
cp ca.crt client.crt client.key $GATEWAY_DIR/

# For Mosquitto (Need CA, Server Cert, Server Key)
# (Need a 'sudo' command because ownership 1883 anda group 1883)
sudo cp ca.crt server.crt server.key $MOSQUITTO_DIR/
sudo chown -R 1883:1883 $MOSQUITTO_DIR

sudo chmod 644 $MOSQUITTO_DIR/ca.crt
sudo chmod 644 $MOSQUITTO_DIR/server.crt
sudo chmod 600 $MOSQUITTO_DIR/server.key

# Cleanup temporary file
rm *.csr *.srl server_ext.cnf ca.key server.key server.crt client.key client.crt ca.crt

echo "Done! Certificates are ready in $GATEWAY_DIR and $MOSQUITTO_DIR"
