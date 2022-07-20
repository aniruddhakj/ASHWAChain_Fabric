/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	handlerContract := new(HandlerContract)
	handlerContract.Info.Version = "0.0.1"
	handlerContract.Info.Description = "My Smart Contract"
	handlerContract.Info.License = new(metadata.LicenseMetadata)
	handlerContract.Info.License.Name = "Apache-2.0"
	handlerContract.Info.Contact = new(metadata.ContactMetadata)
	handlerContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(handlerContract)
	chaincode.Info.Title = "raft chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from HandlerContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
