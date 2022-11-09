/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	// "github.com/hyperledger/fabric/common/flogging"
)

// var spydraLogger = flogging.MustGetLogger("spydraLogger")

type SmartContract struct {
	contractapi.Contract
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		log.Panicf("Error creating spydra chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting spydra chaincode: %s", err.Error())
	}
}
