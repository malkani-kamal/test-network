#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    local MSP_SIGNCERT=$(one_line_pem $6)
    local MSP_KEYSTORE=$(one_line_pem $7)
    local ORDERER_TLS=$(one_line_pem $8)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        -e "s#\${MSP_SIGNCERT}#$MSP_SIGNCERT#" \
        -e "s#\${MSP_KEYSTORE}#$MSP_KEYSTORE#" \
        -e "s#\${ORDERER_TLS}#$ORDERER_TLS#" \
        organizations/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

ORG=1
P0PORT=7051
CAPORT=7054
PEERPEM=organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
CAPEM=organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
MSP_SIGNCERT=organizations/peerOrganizations/org1.example.com/msp/signcerts/cert.pem
MSP_KEYSTORE=$(find organizations/peerOrganizations/org1.example.com/msp/keystore/*)
ORDERER_TLS=organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $MSP_SIGNCERT $MSP_KEYSTORE $ORDERER_TLS)" > organizations/peerOrganizations/org1.example.com/connection-org1.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $MSP_SIGNCERT $MSP_KEYSTORE $ORDERER_TLS)" > organizations/peerOrganizations/org1.example.com/connection-org1.yaml

ORG=2
P0PORT=9051
CAPORT=8054
PEERPEM=organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
CAPEM=organizations/peerOrganizations/org2.example.com/ca/ca.org2.example.com-cert.pem
MSP_SIGNCERT=organizations/peerOrganizations/org2.example.com/msp/signcerts/cert.pem
MSP_KEYSTORE=$(find organizations/peerOrganizations/org2.example.com/msp/keystore/*)

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $MSP_SIGNCERT $MSP_KEYSTORE $ORDERER_TLS)" > organizations/peerOrganizations/org2.example.com/connection-org2.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $MSP_SIGNCERT $MSP_KEYSTORE $ORDERER_TLS)" > organizations/peerOrganizations/org2.example.com/connection-org2.yaml
